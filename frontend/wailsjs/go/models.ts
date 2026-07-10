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
	    id: number;
	    version: string;
	    change_log: string[];
	    create_time: number;
	    create_date: string;
	
	    static createFrom(source: any = {}) {
	        return new VersionInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.version = source["version"];
	        this.change_log = source["change_log"];
	        this.create_time = source["create_time"];
	        this.create_date = source["create_date"];
	    }
	}

}

