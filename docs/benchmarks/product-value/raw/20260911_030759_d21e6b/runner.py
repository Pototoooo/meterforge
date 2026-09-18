#!/usr/bin/env python3
"""Controlled black-box product acceptance: actual local API, synthetic input.
No claim about human task efficiency or authenticated tenant isolation.
"""
import collections, datetime as dt, hashlib, json, math, os, pathlib, random, subprocess, time, urllib.request, urllib.error
ROOT=pathlib.Path(__file__).resolve().parent
BASE=os.environ.get('METERFORGE_ADDRESS','http://127.0.0.1:39888').rstrip('/')
RUN=dt.datetime.now(dt.timezone.utc).strftime('%Y%m%d_%H%M%S')+'_'+os.urandom(3).hex()
OUT=ROOT/'raw'/RUN; OUT.mkdir(parents=True)
START=time.monotonic(); requests=[]; results=[]; registry=[]
COMPOSE=['docker','compose','-p','mf-product-value','-f',str(ROOT/'compose.json')]
def save(name,obj): (OUT/name).write_text(json.dumps(obj,ensure_ascii=False,indent=2)+'\n')
def http(method,path,body=None):
    data=json.dumps(body,separators=(',',':')).encode() if body is not None else None
    start=time.monotonic()
    r=urllib.request.Request(BASE+path,data=data,method=method,headers={'Content-Type':'application/json'})
    try:
        with urllib.request.urlopen(r,timeout=20) as x: status=x.status; text=x.read().decode()
    except urllib.error.HTTPError as x: status=x.code;text=x.read().decode()
    except Exception as x: status=0;text=str(x)
    try: parsed=json.loads(text)
    except ValueError: parsed=text
    record={'n':len(requests)+1,'at':dt.datetime.now(dt.timezone.utc).isoformat(),'elapsed_ms':(time.monotonic()-start)*1000,'method':method,'path':path,'request':body,'status':status,'response':parsed}
    requests.append(record)
    with (OUT/'http.jsonl').open('a') as f:f.write(json.dumps(record,ensure_ascii=False)+'\n')
    return status,parsed

def meter(case, value='$.quantity', groups=None):
    slug='pv_'+RUN+'_'+case
    obj={'slug':slug,'eventType':slug,'aggregation':'SUM','valueProperty':value}
    if groups: obj['groupBy']={g:'$.'+g for g in groups}
    st,data=http('POST','/api/v1/meters',obj)
    if st not in [200,201]: raise RuntimeError('create meter: '+str(data))
    registry.append(slug);return slug

def event(slug,i,quantity=1,subject='a',**extra):
    data={'quantity':quantity,'unit':'token','model':'llm','billing_mode':'platform','price_version':'fixture-v1','job_id':f'job-{i}'};data.update(extra)
    return {'specversion':'1.0','id':slug+'_'+str(i),'source':'product-acceptance','type':slug,'subject':RUN+'_'+subject,'time':(dt.datetime.now(dt.timezone.utc)-dt.timedelta(minutes=2)).isoformat(),'data':data}

def send(events):
    codes=[]
    for i in range(0,len(events),100):
        st,_=http('POST','/api/v1/events',events[i:i+100]);codes.append(st)
    return codes

def query(slug,extra=''):
    return http('GET','/api/v1/meters/'+slug+'/query'+extra)
def total(data):
    if not isinstance(data,dict) or not isinstance(data.get('data'),list): return None
    return sum(float(r['value']) for r in data['data'])
def wait_total(slug,expected,timeout=45,extra=''):
    t=time.monotonic();last=None;status=None
    while time.monotonic()-t<timeout:
        status,last=query(slug,extra)
        if status==200 and total(last) is not None and abs(total(last)-expected)<1e-6:
            return True,last,round(time.monotonic()-t,3)
        time.sleep(1)
    return False,last,round(time.monotonic()-t,3)
def result(case,passed,**facts):
    r={'case':case,'passed':bool(passed),**facts};results.append(r);save('results.json',results);print(json.dumps(r,ensure_ascii=False),flush=True)
def main():
    st,_=http('GET','/api/v1/meters')
    if st!=200: raise RuntimeError('MeterForge unavailable')
    # Five independent model meters; never sum seconds and tokens.
    for kind in ['llm','embedding','rerank','vlm','asr']:
        slug=meter(kind,groups=['model','billing_mode','price_version']); t=time.monotonic()
        inputs=[event(slug,i,quantity=(i%17)+1,subject='a' if i%2==0 else 'b',model=kind,unit='second' if kind=='asr' else 'token',billing_mode='platform' if i%3 else 'byok') for i in range(1000)]
        random.Random(20260911).shuffle(inputs)
        for i,e in enumerate(inputs):
            if i%5==0: e['time']=(dt.datetime.now(dt.timezone.utc)-dt.timedelta(days=2,minutes=i%60)).isoformat()
        save(kind+'-oracle.json',{'unique_events':1000,'expected_total':sum(e['data']['quantity'] for e in inputs),'expected_by_subject':{RUN+'_'+s:sum(e['data']['quantity'] for e in inputs if e['subject']==RUN+'_'+s) for s in ['a','b']}})
        codes=send(inputs);expected=sum(e['data']['quantity'] for e in inputs)
        ok,observed,seconds=wait_total(slug,expected)
        result(kind+'_unique',ok and all(c==204 for c in codes),unique_events=1000,expected=expected,observed=total(observed),batch_statuses=codes,visible_after_last_send_s=seconds,total_machine_task_s=round(time.monotonic()-t,3))
        # Distribute repeated events across both subjects and irregular quantities.
        duplicate=inputs[:250];codes=send(duplicate);time.sleep(4)
        ok,observed,seconds=wait_total(slug,expected,15)
        result(kind+'_duplicate',ok and all(c==204 for c in codes),duplicate_events=250,expected_after_replay=expected,observed=total(observed),batch_statuses=codes,observation_delay_s=4,qualification='Finite settled-window aggregate check, not lifetime exactly-once proof')
        # Exact subject filtering is attribution, not tenant authorization.
        for subject in ['a','b']:
            exp=sum(e['data']['quantity'] for e in inputs if e['subject']==RUN+'_'+subject)
            ok,obs,seconds=wait_total(slug,exp,5,'?subject='+RUN+'_'+subject)
            result(kind+'_subject_'+subject,ok,expected=exp,observed=total(obs),scope='subject-filter attribution only')
        sample=inputs[0]
        st,events=http('GET','/api/v1/events?limit=10&id='+sample['id'])
        raw=events if isinstance(events,list) else events.get('data',[]) if isinstance(events,dict) else []
        matches=[r for r in raw if isinstance(r,dict) and r.get('event',{}).get('id')==sample['id']]
        complete=any(all(k in r['event'].get('data',{}) for k in ['model','unit','job_id','billing_mode','price_version','quantity']) and r['event']['subject']==sample['subject'] for r in matches)
        result(kind+'_trace',st==200 and complete,sample_size=1,exact_event_matches=len(matches),http_status=st)
    # Envelopes vs semantic unit validation: report actual layers separately.
    slug=meter('invalid')
    for label,change in [('missing_id',lambda e:e.pop('id')),('missing_source',lambda e:e.pop('source')),('invalid_specversion',lambda e:e.update(specversion='invalid')),('invalid_time',lambda e:e.update(time='not-a-time'))]:
        e=event(slug,label);change(e);st,body=http('POST','/api/v1/events',e)
        result(label,400<=st<500,http_status=st,expected='4xx synchronous rejection')
    # Negative and wrong-unit events are business-invalid for this fixture, not
    # necessarily invalid CloudEvents. Test producer-contract assumptions honestly.
    slug=meter('semantics')
    for label,quantity,unit in [('wrong_unit',10,'banana'),('negative',-2,'token'),('missing_quantity',None,'token'),('nonnumeric','ten','token')]:
        e=event(slug,label,quantity,unit=unit)
        if label=='missing_quantity':e['data'].pop('quantity')
        st,body=http('POST','/api/v1/events',e)
        result(label+'_rejection',400<=st<500,http_status=st,expected='business-invalid input rejected',scope='producer semantic contract probe, not a CloudEvents standards assertion')
    time.sleep(4);st,obs=query(slug);save('semantic-query.json',obs)
    st,invalids=http('GET','/api/v1/events?limit=100&subject='+RUN+'_a');save('semantic-events.json',invalids)
    # Scope spoofing probe: default quickstart is a trusted producer API.
    slug=meter('subject_claim')
    e=event(slug,0,999,subject='other_workspace');st,_=http('POST','/api/v1/events',e)
    ok,obs,_=wait_total(slug,999,10)
    result('subject_claim_rejected',st in [401,403],http_status=st,observed_quantity=total(obs),scope='unauthenticated local default namespace; cannot claim authenticated tenant isolation')
    # Controlled consumer outage in OUR disposable project only.
    slug=meter('recovery')
    baseline=[event(slug,i) for i in range(100)];codes=send(baseline);ok,obs,_=wait_total(slug,100)
    if not ok:raise RuntimeError('recovery baseline not ready')
    paused=False
    try:
        p=subprocess.run(COMPOSE+['pause','sink-worker'],capture_output=True,text=True);save('pause.json',{'command':COMPOSE+['pause','sink-worker'],'exit':p.returncode,'stdout':p.stdout,'stderr':p.stderr});p.check_returncode();paused=True
        outage=[event(slug,100+i) for i in range(500)];codes=send(outage);time.sleep(5);st,obs=query(slug)
        result('consumer_paused',all(c==204 for c in codes) and total(obs)==100,accepted_events=500,batch_statuses=codes,visible_quantity_during_pause=total(obs),pause_hold_after_send_s=5)
    finally:
        if paused:
            p=subprocess.run(COMPOSE+['unpause','sink-worker'],capture_output=True,text=True);save('unpause.json',{'exit':p.returncode,'stdout':p.stdout,'stderr':p.stderr});p.check_returncode()
    ok,obs,seconds=wait_total(slug,600,60)
    result('consumer_recovered',ok,expected=600,observed=total(obs),recovery_to_visible_s=seconds,backlog_events=500)
    # Replay backlog after recovery, verify stable result after settling again.
    codes=send(outage);time.sleep(4);ok,obs,seconds=wait_total(slug,600,15)
    result('recovery_replay',ok and all(c==204 for c in codes),replayed_events=500,expected=600,observed=total(obs))
    # Late-arriving facts after the initial window was queried.
    slug=meter('late')
    a=[event(slug,i,2) for i in range(200)];codes=send(a);ok,obs,_=wait_total(slug,400)
    b=[event(slug,200+i,3) for i in range(200)]
    for e in b:e['time']=(dt.datetime.now(dt.timezone.utc)-dt.timedelta(days=7)).isoformat()
    codes=send(b);ok,obs,seconds=wait_total(slug,1000)
    result('late_arrival',ok and all(c==204 for c in codes),initial_events=200,late_events=200,late_age_days=7,expected=1000,observed=total(obs),visible_after_last_send_s=seconds)
    # Clean empty meter control makes stale data and missing filter bugs visible.
    slug=meter('empty');st,obs=query(slug)
    result('empty_meter',st==200 and total(obs)==0,observed=total(obs))

try:
    main();status='completed'
except Exception as e:
    status='execution_error';save('error.json',{'error':repr(e)});print(repr(e),flush=True)
finally:
    summary={'run':RUN,'status':status,'base_url':BASE,'seed':20260911,'results':results,'meters':registry,'http_requests':len(requests),'duration_s':round(time.monotonic()-START,3),'human_participants':0,'actual_payments':0,'source':'actual MeterForge local image / synthetic CloudEvents; not live LoreLattice traffic'}
    save('SUMMARY.json',summary);(ROOT/'LATEST_RUN').write_text(RUN+'\n')
    print('RESULT_DIR='+str(OUT),flush=True)
if status!='completed':raise SystemExit(2)
# Product gaps are recorded as failed cases, not hidden by a nonzero shell exit.
