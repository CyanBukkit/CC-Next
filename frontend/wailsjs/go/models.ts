export namespace config {
	
	export class Settings {
	    passphrase_length: number;
	    stuck_timeout_seconds: number;
	    auto_continue_mode: boolean;
	    voice_provider: string;
	    claude_model: string;
	    dark_mode: boolean;
	    max_history_messages: number;
	    passphrase_template: string;
	    work_dir: string;
	    provider_auth_token: string;
	    provider_base_url: string;
	    provider_model: string;
	
	    static createFrom(source: any = {}) {
	        return new Settings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.passphrase_length = source["passphrase_length"];
	        this.stuck_timeout_seconds = source["stuck_timeout_seconds"];
	        this.auto_continue_mode = source["auto_continue_mode"];
	        this.voice_provider = source["voice_provider"];
	        this.claude_model = source["claude_model"];
	        this.dark_mode = source["dark_mode"];
	        this.max_history_messages = source["max_history_messages"];
	        this.passphrase_template = source["passphrase_template"];
	        this.work_dir = source["work_dir"];
	        this.provider_auth_token = source["provider_auth_token"];
	        this.provider_base_url = source["provider_base_url"];
	        this.provider_model = source["provider_model"];
	    }
	}

}

export namespace main {
	
	export class SessionInfo {
	    id: string;
	    title: string;
	    created_at: string;
	    updated_at: string;
	
	    static createFrom(source: any = {}) {
	        return new SessionInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.created_at = source["created_at"];
	        this.updated_at = source["updated_at"];
	    }
	}

}

export namespace message {
	
	export class Message {
	    role: string;
	    content: string;
	
	    static createFrom(source: any = {}) {
	        return new Message(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.role = source["role"];
	        this.content = source["content"];
	    }
	}

}

