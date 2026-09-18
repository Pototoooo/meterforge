#!/usr/bin/env python3
"""Recompute denominators from recorded wire inputs, not precomputed claims."""
import collections, csv, hashlib, json, pathlib, statistics, urllib.request, sys
ROOT=pathlib.Path(__file__).resolve().parent
runs=[]; metrics=[]; audits=[]
offline='--offline' in sys.argv
recorded={(a['run'],a['kind']):a for a in json.loads((ROOT/'FINAL_READBACK.json').read_text())} if offline else {}
for path in sorted((ROOT/'raw').glob('*/SUMMARY.json')):
    s=json.loads(path.read_text())
    if s.get('cohort')!='mapped_customer_v2':continue
    assert s['status']=='completed',path
    wire=[json.loads(x) for x in (path.parent/'http.jsonl').read_text().splitlines()]
    look={r['case']:r for r in s['results']}
    for kind in ['llm','embedding','rerank','vlm','asr']:
        typ='pv_'+s['run']+'_'+kind
        submissions=[]
        for call in wire:
            if call['method']=='POST' and call['path']=='/api/v1/events':
                events=call['request'] if isinstance(call['request'],list) else [call['request']]
                selected=[e for e in events if e.get('type')==typ]
                if selected:assert call['status']==204
                submissions+=selected
        unique={}
        for e in submissions:
            key=(e['source'],e['id'])
            if key in unique:assert unique[key]==e
            unique[key]=e
        assert len(unique)==1000 and len(submissions)==1250
        groups=collections.defaultdict(int)
        for e in unique.values():
            groups[(e['subject'],e['data']['model'],e['data']['billing_mode'],e['data']['price_version'])]+=e['data']['quantity']
        expected=sum(e['data']['quantity'] for e in unique.values())
        expected_cost=sum(e['data']['cost_microusd'] for e in unique.values())
        actual=look[kind+'_unique']['observed']; duplicate_actual=look[kind+'_duplicate']['observed']
        assert actual==expected==duplicate_actual
        assert look[kind+'_cost_metadata']['observed_microusd']==expected_cost
        assert look[kind+'_dimensions']['matched_groups']==len(groups)==8
        grouped_calls=[c for c in wire if c['method']=='GET' and c['path'].startswith('/api/v1/meters/'+typ+'/query?groupBy=')]
        assert len(grouped_calls)==1 and grouped_calls[0]['status']==200
        wire_groups=collections.defaultdict(float)
        for row in grouped_calls[0]['response']['data']:
            d=row['groupBy']; wire_groups[(row['subject'],d['model'],d['billing_mode'],d['price_version'])]+=row['value']
        assert dict(wire_groups)==dict(groups)

        # A late final read bounds delayed duplicate-processing uncertainty.
        url=s['base_url']+'/api/v1/meters/'+typ+'/query'
        if offline:raw=recorded[(s['run'],kind)]['final_response']
        else:
            with urllib.request.urlopen(url,timeout=15) as r:raw=json.load(r)
        late_total=sum(row['value'] for row in raw['data']);assert late_total==expected
        audits.append({'run':s['run'],'kind':kind,'url':url,'final_response':raw,'expected':expected,'matches':late_total==expected})
        metrics.append({'run':s['run'],'model':kind,'unique_events':len(unique),'retries':len(submissions)-len(unique),'expected_quantity':expected,'actual_quantity':actual,'after_retries':duplicate_actual,'final_read':late_total,'groups':len(groups),'cost_microusd':expected_cost,'batch_visibility_s':look[kind+'_unique']['visible_after_last_send_s'],'machine_task_s':look[kind+'_unique']['total_machine_task_s']})
    runs.append(s)
assert len(runs)==3, 'Expected exactly three qualification runs'
(ROOT/'FINAL_READBACK.json').write_text(json.dumps(audits,ensure_ascii=False,indent=2)+'\n')
with (ROOT/'METRICS.csv').open('w') as f:
    w=csv.DictWriter(f,fieldnames=list(metrics[0]));w.writeheader();w.writerows(metrics)
latency=[r['batch_visibility_s'] for r in metrics]
recovery=[next(r for r in s['results'] if r['case']=='consumer_recovered')['recovery_to_visible_s'] for s in runs]
rollup={'runs':[s['run'] for s in runs],'model_run_checks':len(metrics),'unique_main_events':sum(r['unique_events'] for r in metrics),'main_retries':sum(r['retries'] for r in metrics),'group_checks':sum(r['groups'] for r in metrics),'trace_samples':15,'cost_metadata_checks':15,'quantity_difference':sum(abs(r['actual_quantity']-r['expected_quantity']) for r in metrics),'batch_visibility_median_s':statistics.median(latency),'batch_visibility_min_s':min(latency),'batch_visibility_max_s':max(latency),'recovery_seconds':recovery,'recovery_backlog_events':1500,'post_recovery_retries':1500,'late_events':600,'envelope_invalid_tests':12,'business_invalid_tests':12,'subject_claim_probes':3,'human_participants':0,'actual_payments':0,'failed_cases':[{'run':s['run'],'cases':[r['case'] for r in s['results'] if not r['passed']]} for s in runs]}
(ROOT/'SUMMARY.json').write_text(json.dumps(rollup,ensure_ascii=False,indent=2)+'\n')
print(json.dumps(rollup,ensure_ascii=False,indent=2))
