declare module "*.module.css" {
	const classes: any;
	export default classes;
}

declare module "*.css" {
	const content: string;
	export default content;
}
