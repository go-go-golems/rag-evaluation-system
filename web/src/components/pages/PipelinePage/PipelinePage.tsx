import { PipelineOverview, type PipelineOverviewProps } from '../../pipeline/PipelineOverview';

export type PipelinePageProps = PipelineOverviewProps;

export function PipelinePage(props: PipelinePageProps) {
  return <PipelineOverview {...props} />;
}
