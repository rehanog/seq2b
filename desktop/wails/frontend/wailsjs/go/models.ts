export namespace main {
	
	export class BacklinkData {
	    sourcePage: string;
	    blockIds: string[];
	    count: number;
	
	    static createFrom(source: any = {}) {
	        return new BacklinkData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.sourcePage = source["sourcePage"];
	        this.blockIds = source["blockIds"];
	        this.count = source["count"];
	    }
	}
	export class SegmentData {
	    type: string;
	    content: string;
	    target?: string;
	    alt?: string;
	
	    static createFrom(source: any = {}) {
	        return new SegmentData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.content = source["content"];
	        this.target = source["target"];
	        this.alt = source["alt"];
	    }
	}
	export class BlockData {
	    id: string;
	    content: string;
	    htmlContent: string;
	    segments: SegmentData[];
	    depth: number;
	    children: BlockData[];
	    todoState: string;
	    checkboxState: string;
	    priority: string;
	
	    static createFrom(source: any = {}) {
	        return new BlockData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.content = source["content"];
	        this.htmlContent = source["htmlContent"];
	        this.segments = this.convertValues(source["segments"], SegmentData);
	        this.depth = source["depth"];
	        this.children = this.convertValues(source["children"], BlockData);
	        this.todoState = source["todoState"];
	        this.checkboxState = source["checkboxState"];
	        this.priority = source["priority"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class PageData {
	    name: string;
	    title: string;
	    blocks: BlockData[];
	    backlinks: BacklinkData[];
	
	    static createFrom(source: any = {}) {
	        return new PageData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.title = source["title"];
	        this.blocks = this.convertValues(source["blocks"], BlockData);
	        this.backlinks = this.convertValues(source["backlinks"], BacklinkData);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

