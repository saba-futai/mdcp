export namespace models {
	
	export class RandomOptions {
	    count: number;
	    minLevel: number;
	    maxLevel: number;
	    onlyFavorites: boolean;
	
	    static createFrom(source: any = {}) {
	        return new RandomOptions(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.count = source["count"];
	        this.minLevel = source["minLevel"];
	        this.maxLevel = source["maxLevel"];
	        this.onlyFavorites = source["onlyFavorites"];
	    }
	}
	export class Song {
	    id: number;
	    title: string;
	    artist: string;
	    level: number;
	    color: string;
	    isFavorite: boolean;
	    isBlacklisted: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Song(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.artist = source["artist"];
	        this.level = source["level"];
	        this.color = source["color"];
	        this.isFavorite = source["isFavorite"];
	        this.isBlacklisted = source["isBlacklisted"];
	    }
	}

}

