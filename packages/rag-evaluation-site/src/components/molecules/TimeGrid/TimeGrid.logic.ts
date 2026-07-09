export interface TimeGridLogicBlock {
	id: string;
	startISO: string;
	endISO: string;
}

export interface TimeParts {
	date: string;
	minutes: number;
}

export interface PackedTimeGridBlock<Block extends TimeGridLogicBlock = TimeGridLogicBlock> {
	block: Block;
	topPct: number;
	heightPct: number;
	lane: number;
	lanes: number;
}

export function timeParts(iso: string): TimeParts {
	const date = iso.slice(0, 10);
	const hh = Number(iso.slice(11, 13));
	const mm = Number(iso.slice(14, 16));
	const minutes = (Number.isFinite(hh) ? hh : 0) * 60 + (Number.isFinite(mm) ? mm : 0);
	return { date, minutes };
}

/**
 * Pack overlapping timed blocks into side-by-side lanes. Blocks are grouped into
 * overlap clusters; within a cluster each block gets a lane and every member
 * shares the cluster's lane count so widths line up.
 */
export function packTimeGridColumn<Block extends TimeGridLogicBlock>(
	blocks: Block[],
	rangeStart: number,
	rangeMinutes: number,
): PackedTimeGridBlock<Block>[] {
	const timed = blocks
		.map((block) => {
			const start = timeParts(block.startISO).minutes;
			const end = Math.max(start + 15, timeParts(block.endISO).minutes);
			return { block, start, end };
		})
		.filter((b) => b.end > rangeStart && b.start < rangeStart + rangeMinutes)
		.sort((a, b) => a.start - b.start || a.end - b.end);

	const packed: PackedTimeGridBlock<Block>[] = [];
	let cluster: Array<(typeof timed)[number] & { lane: number }> = [];
	let clusterEnd = -1;
	const laneEnds: number[] = [];

	const flush = () => {
		const lanes = laneEnds.length || 1;
		for (const item of cluster) {
			const start = Math.max(item.start, rangeStart);
			const end = Math.min(item.end, rangeStart + rangeMinutes);
			packed.push({
				block: item.block,
				topPct: ((start - rangeStart) / rangeMinutes) * 100,
				heightPct: ((end - start) / rangeMinutes) * 100,
				lane: item.lane,
				lanes,
			});
		}
		cluster = [];
		laneEnds.length = 0;
		clusterEnd = -1;
	};

	for (const item of timed) {
		if (cluster.length > 0 && item.start >= clusterEnd) flush();
		let lane = laneEnds.findIndex((end) => end <= item.start);
		if (lane === -1) {
			lane = laneEnds.length;
			laneEnds.push(item.end);
		} else {
			laneEnds[lane] = item.end;
		}
		cluster.push({ ...item, lane });
		clusterEnd = Math.max(clusterEnd, item.end);
	}
	if (cluster.length > 0) flush();
	return packed;
}
