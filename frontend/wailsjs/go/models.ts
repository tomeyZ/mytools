export namespace handler {
	
	export class CryptoResponse {
	    success: boolean;
	    result: string;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new CryptoResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.result = source["result"];
	        this.message = source["message"];
	    }
	}
	export class IPInfo {
	    status: string;
	    message?: string;
	    country: string;
	    countryCode: string;
	    region: string;
	    regionName: string;
	    city: string;
	    zip: string;
	    lat: number;
	    lon: number;
	    timezone: string;
	    isp: string;
	    org: string;
	    as: string;
	    query: string;
	
	    static createFrom(source: any = {}) {
	        return new IPInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.status = source["status"];
	        this.message = source["message"];
	        this.country = source["country"];
	        this.countryCode = source["countryCode"];
	        this.region = source["region"];
	        this.regionName = source["regionName"];
	        this.city = source["city"];
	        this.zip = source["zip"];
	        this.lat = source["lat"];
	        this.lon = source["lon"];
	        this.timezone = source["timezone"];
	        this.isp = source["isp"];
	        this.org = source["org"];
	        this.as = source["as"];
	        this.query = source["query"];
	    }
	}
	export class RSAKeyPair {
	    success: boolean;
	    publicKey: string;
	    privateKey: string;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new RSAKeyPair(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.publicKey = source["publicKey"];
	        this.privateKey = source["privateKey"];
	        this.message = source["message"];
	    }
	}
	export class VersionInfo {
	    version: string;
	    change_log: string[];
	    create_date: string;
	    download_url: string;
	    release_url: string;
	
	    static createFrom(source: any = {}) {
	        return new VersionInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.version = source["version"];
	        this.change_log = source["change_log"];
	        this.create_date = source["create_date"];
	        this.download_url = source["download_url"];
	        this.release_url = source["release_url"];
	    }
	}
	export class UpdateResult {
	    status: string;
	    message: string;
	    latest?: VersionInfo;
	
	    static createFrom(source: any = {}) {
	        return new UpdateResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.status = source["status"];
	        this.message = source["message"];
	        this.latest = this.convertValues(source["latest"], VersionInfo);
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

