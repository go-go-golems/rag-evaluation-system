#!/usr/bin/env python3
import argparse, csv, json
from datetime import datetime
from pathlib import Path
import matplotlib.pyplot as plt

p=argparse.ArgumentParser()
p.add_argument("evidence")
p.add_argument("--output",required=True)
a=p.parse_args()
data=json.loads(Path(a.evidence).read_text())
out=Path(a.output); out.mkdir(parents=True,exist_ok=True)
cells=data["cells"]
colors={1:"#1f77b4",2:"#ff7f0e",4:"#2ca02c"}

def series(field):
    result={}
    for cell in cells:
        c=cell["cell"]["concurrency"]
        result.setdefault(c,[]).append((cell["cell"]["chunksPerRequest"],cell[field]))
    return {k:sorted(v) for k,v in result.items()}

def plot_metric(field,ylabel,name,unavailable_if_zero=False):
    fig,ax=plt.subplots(figsize=(8,5))
    unavailable = unavailable_if_zero and all(cell[field] == 0 for cell in cells)
    if unavailable:
        ax.text(.5,.5,"Data unavailable\nFixture provider reports no token/cost usage",ha="center",va="center",transform=ax.transAxes,fontsize=14,color="#555555")
    for concurrency,values in sorted(series(field).items()):
        if not unavailable:
            ax.plot([x for x,_ in values],[y for _,y in values],marker="o",color=colors[concurrency],label=f"concurrency={concurrency}")
    ax.set(xlabel="Chunks per LLM request",ylabel=ylabel,title=f"Workflow V3 fixture control — {ylabel}")
    ax.set_xticks([1,2,4,8]); ax.grid(True,alpha=.25)
    if not unavailable: ax.legend()
    fig.tight_layout()
    for ext in ("svg","png"): fig.savefig(out/f"{name}.{ext}",dpi=160)
    plt.close(fig)

plot_metric("chunksPerSecond","Chunks / second","chunks-throughput")
plot_metric("requestsPerSecond","Requests / second","request-throughput")
for cell in cells:
    elapsed = cell["makespanMicros"] / 1_000_000
    usage = cell.get("usage", {})
    cell["tokensPerSecond"] = (usage.get("input_tokens", 0) + usage.get("output_tokens", 0)) / elapsed if elapsed else 0
    cell["costMicrounitsPerChunk"] = usage.get("cost_microunits", 0) / cell["chunks"]
plot_metric("tokensPerSecond","Tokens / second (not reported by fixture)","token-rate",True)
plot_metric("costMicrounitsPerChunk","Cost, microunits / chunk (not reported by fixture)","cost-efficiency",True)
for cell in cells:
    cell["overlapMillis"] = cell.get("overlapMicros", 0) / 1000
if all(cell["overlapMillis"] == 0 for cell in cells):
    fig,ax=plt.subplots(figsize=(8,5)); ax.axhline(0,color="#333333",linewidth=2)
    ax.text(.5,.55,"Observed zero in all fixture cells\nAll concurrency series coincide at 0 ms\nNot indicative of real-provider performance",ha="center",va="center",transform=ax.transAxes,fontsize=13,color="#555555")
    ax.set(xlabel="Chunks per LLM request",ylabel="Generation / embedding overlap (ms)",title="Workflow V3 fixture control — generation / embedding overlap")
    ax.set_xticks([1,2,4,8]); ax.set_ylim(0,.05); ax.grid(True,alpha=.25); fig.tight_layout()
    for ext in ("svg","png"): fig.savefig(out/f"generation-embedding-overlap.{ext}",dpi=160)
    plt.close(fig)
else:
    plot_metric("overlapMillis","Generation / embedding overlap (ms)","generation-embedding-overlap")
for cell in cells: cell["makespanMillis"] = cell["makespanMicros"] / 1000
plot_metric("makespanMillis","Cell makespan (ms)","makespan")

fig,ax=plt.subplots(figsize=(9,5))
offsets={1:-0.08,2:0,4:0.08}
for cell in cells:
    batch=cell["cell"]["chunksPerRequest"]; concurrency=cell["cell"]["concurrency"]
    values=sorted(cell["providerMicros"])
    median=values[len(values)//2]
    ax.scatter(batch+offsets[concurrency],median,label=f"concurrency={concurrency}" if batch==1 else None,s=50,color=colors[concurrency])
ax.set(xlabel="Chunks per LLM request",ylabel="Median provider span (µs)",title="Workflow V3 fixture control — provider wall time")
ax.set_xticks([1,2,4,8]); ax.grid(True,alpha=.25)
handles,labels=ax.get_legend_handles_labels(); ax.legend(handles,labels)
fig.tight_layout()
for ext in ("svg","png"): fig.savefig(out/f"provider-latency.{ext}",dpi=160)
plt.close(fig)

fig,ax=plt.subplots(figsize=(10,6))
rows=sorted([cell for cell in cells if cell["cell"]["concurrency"]==4],key=lambda x:x["cell"]["chunksPerRequest"])
batch_colors={1:"#1f77b4",2:"#ff7f0e",4:"#2ca02c",8:"#d62728"}; peak=4
for cell in rows:
    batch=cell["cell"]["chunksPerRequest"]
    all_attempts=cell["attempts"]+cell["embeddingAttempts"]
    origin=min(datetime.fromisoformat(a["startedAt"].replace("Z","+00:00")) for a in all_attempts)
    for attempts,phase,style in ((cell["attempts"],"generation","-"),(cell["embeddingAttempts"],"embedding",":")):
        events=[]
        for attempt in attempts:
            start=datetime.fromisoformat(attempt["startedAt"].replace("Z","+00:00")); end=datetime.fromisoformat(attempt["finishedAt"].replace("Z","+00:00"))
            events.extend([((start-origin).total_seconds(),1),((end-origin).total_seconds(),-1)])
        events.sort(key=lambda item:(item[0],item[1])); x=[0.0]; y=[0]; active=0
        for stamp,delta in events: x.extend([stamp,stamp]); y.extend([active,active+delta]); active+=delta
        peak=max(peak,max(y,default=0)); ax.plot(x,y,color=batch_colors[batch],linestyle=style,label=f"batch={batch} {phase}")
ax.axhline(4,color="#333333",linestyle="--",linewidth=1,label="Umans generation hard cap=4")
ax.set(xlabel="Seconds since cell admission",ylabel="Active attempts",title="Workflow V3 fixture control — generation and embedding activity")
ax.set_xlim(left=0); ax.set_yticks(range(0,peak+1)); ax.grid(True,alpha=.25); ax.legend(ncol=2,fontsize=9); fig.tight_layout()
for ext in ("svg","png"): fig.savefig(out/f"request-timeline.{ext}",dpi=160)
plt.close(fig)

for svg in out.glob("*.svg"):
    svg.write_text("\n".join(line.rstrip() for line in svg.read_text().splitlines())+"\n")
summary={"schemaVersion":"rag-ttc-v3-sweep-graph-manifest/v1","evidencePlanDigest":data["plan"]["digest"],"graphs":sorted(x.name for x in out.iterdir() if x.suffix in {".svg",".png"})}
(out/"manifest.json").write_text(json.dumps(summary,indent=2)+"\n")
print(f"rendered={len(summary['graphs'])} output={out}")
