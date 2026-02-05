export namespace wailsapp {
	
	export class UIState {
	    IsActive: boolean;
	    ZoomLevel: number;
	    ScopeSize: number;
	    BorderEnabled: boolean;
	    BorderColor: number;
	    BorderThickness: number;
	    FollowCursor: boolean;
	    HotkeyMode: string;
	    Hotkey: number;
	    RenderBackend: string;
	
	    static createFrom(source: any = {}) {
	        return new UIState(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.IsActive = source["IsActive"];
	        this.ZoomLevel = source["ZoomLevel"];
	        this.ScopeSize = source["ScopeSize"];
	        this.BorderEnabled = source["BorderEnabled"];
	        this.BorderColor = source["BorderColor"];
	        this.BorderThickness = source["BorderThickness"];
	        this.FollowCursor = source["FollowCursor"];
	        this.HotkeyMode = source["HotkeyMode"];
	        this.Hotkey = source["Hotkey"];
	        this.RenderBackend = source["RenderBackend"];
	    }
	}

}

