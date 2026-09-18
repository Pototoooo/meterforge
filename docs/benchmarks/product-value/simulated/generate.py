#!/usr/bin/env python3
"""Illustrative, explicitly synthetic task rows. No real participant observations."""
import csv,json,statistics
from pathlib import Path
ROOT=Path(__file__).resolve().parent
rows=[]
# Outcomes and durations below are authored scenario assumptions, not inferred from API benchmarks.
for i in range(24):
    old_ok=i not in {2,5,8,11,14,20,23}
    new_ok=i not in {11,23}
    for variant,ok,helped,minutes in [
        ('baseline',old_ok,i%3==0 or not old_ok,12+(i*7)%17),
        ('proposed',new_ok,i in {5,11,17,23},8+(i*5)%11)]:
        rows.append(dict(data_origin='synthetic_assumption',study='onboarding',persona_id=f'SYN-D{i//3+1:02}',task_id=['llm_meter','asr_meter','customer_attribution'][i%3],variant=variant,success=int(ok),assisted=int(helped),recorded_minutes=minutes if ok else 30,limit_minutes=30,critical_error=int(not ok)))
for i in range(36):
    old_ok=i not in {0,3,4,7,10,13,15,18,21,24,27,30,33}
    new_ok=i not in {3,10,18,27,33}
    for variant,ok,helped,minutes in [
        ('baseline',old_ok,i%2==0,6+(i*5)%13),
        ('proposed',new_ok,i in {3,10,18,25,27,33},4+(i*3)%9)]:
        rows.append(dict(data_origin='synthetic_assumption',study='operations',persona_id=f'SYN-O{i//6+1:02}',task_id=['missing_customer','invalid_unit','negative_usage','price_version','retry_duplicate','consumer_pause'][i%6],variant=variant,success=int(ok),assisted=int(helped),recorded_minutes=minutes if ok else 20,limit_minutes=20,critical_error=int(i in ({3,15,27} if variant=='baseline' else {27}))))
for i in range(40):
    for variant,correct in [('baseline',i%3!=0),('proposed',i not in {2,9,16,23,30})]:
        rows.append(dict(data_origin='synthetic_assumption',study='understanding',persona_id=f'SYN-U{i//5+1:02}',task_id=['who_pays','cost_explanation','credit_timing','byok','failure_state'][i%5],variant=variant,success=int(correct),assisted=0,recorded_minutes='',limit_minutes='',critical_error=0))
with (ROOT/'SYNTHETIC_TASKS.csv').open('w') as f:
    w=csv.DictWriter(f,fieldnames=list(rows[0]));w.writeheader();w.writerows(rows)
metrics={}
for study in ['onboarding','operations','understanding']:
    metrics[study]={}
    for variant in ['baseline','proposed']:
        r=[x for x in rows if x['study']==study and x['variant']==variant]
        m={'n':len(r),'success':sum(x['success'] for x in r),'success_rate':sum(x['success'] for x in r)/len(r),'assisted':sum(x['assisted'] for x in r),'critical_errors':sum(x['critical_error'] for x in r)}
        if study!='understanding':
            m['independent_success']=sum(x['success'] and not x['assisted'] for x in r)
            m['capped_time_median_minutes']=statistics.median(x['recorded_minutes'] for x in r)
        metrics[study][variant]=m
    b,p=metrics[study]['baseline'],metrics[study]['proposed']
    metrics[study]['delta_pp']=(p['success_rate']-b['success_rate'])*100
    if study!='understanding':metrics[study]['capped_median_reduction_pct']=(1-p['capped_time_median_minutes']/b['capped_time_median_minutes'])*100
(ROOT/'SIMULATED_SUMMARY.json').write_text(json.dumps({'data_origin':'synthetic_assumption','actual_participants':0,'actual_payments':0,'metrics':metrics},ensure_ascii=False,indent=2)+'\n')
print(json.dumps(metrics,ensure_ascii=False,indent=2))
