import type { HTMLAttributes } from "react";
import type { ContextCourse } from "../../../context";
import { Button } from "../../atoms";
import { Caption } from "../../foundation";
import { SectionBlock, SplitPane, Stack } from "../../layout";
import { CheckList, KeyValueStrip, PersonSummary, StepList } from "../../molecules";
import styles from "./CourseLessonPanel.module.css";

export interface CourseLessonPanelProps extends Omit<HTMLAttributes<HTMLDivElement>, "title"> {
	course: ContextCourse;
	activeAgendaItemId?: string;
	onAgendaItemSelect?: (itemId: string) => void;
	onPrimaryCta?: () => void;
	onSecondaryCta?: () => void;
}

export function CourseLessonPanel({
	course,
	activeAgendaItemId,
	onAgendaItemSelect,
	onPrimaryCta,
	onSecondaryCta,
	className,
	...rest
}: CourseLessonPanelProps) {
	const labels = {
		outcomes: "What you'll leave with",
		agenda: "Agenda · 2 hours",
		instructor: "Instructor",
		when: "When",
		where: "Where",
		format: "Format",
		primaryCta: "Reserve a seat",
		secondaryCta: "Preview the deck",
		...course.labels,
	};
	const agendaItems = course.agenda.map((item) => ({
		id: item.id,
		index: item.number,
		title: item.title,
		description: item.description,
		meta: item.duration,
	}));

	return (
		<article
			className={[styles.root, className ?? ""].filter(Boolean).join(" ")}
			data-rag-organism="CourseLessonPanel"
			{...rest}
		>
			<SectionBlock density="spacious" className={styles.hero}>
				{course.kicker && <div className={styles.eyebrow}>{course.kicker}</div>}
				<h1 className={styles.title}>{course.title}</h1>
				<p className={styles.lead}>{course.tagline}</p>
				<div className={styles.ctaRow}>
					<Button variant="primary" onClick={onPrimaryCta}>
						{labels.primaryCta}
					</Button>
					<Button onClick={onSecondaryCta}>{labels.secondaryCta}</Button>
					{course.price && <Caption>{course.price}</Caption>}
				</div>
				<KeyValueStrip
					className={styles.metaStrip}
					items={[
						{ key: labels.when, value: course.when ?? "TBD" },
						{ key: labels.where, value: course.where ?? "TBD" },
						{ key: labels.format, value: course.format ?? "TBD" },
					]}
				/>
			</SectionBlock>

			<SplitPane
				ratio="course"
				left={
					<SectionBlock label={labels.outcomes} density="spacious">
						<Stack gap="lg">
							<CheckList items={course.outcomes} />
							{course.instructor && (
								<div className={styles.instructorBlock}>
									<div className={styles.sectionLabel}>{labels.instructor}</div>
									<PersonSummary
										name={course.instructor.name}
										subtitle={course.instructor.role}
										bio={course.instructor.bio}
									/>
								</div>
							)}
						</Stack>
					</SectionBlock>
				}
				right={
					<SectionBlock label={labels.agenda} density="spacious">
						<StepList items={agendaItems} />
					</SectionBlock>
				}
			/>

			{course.blurb && (
				<SectionBlock density="spacious">
					<p className={styles.blurb}>{course.blurb}</p>
				</SectionBlock>
			)}
		</article>
	);
}
