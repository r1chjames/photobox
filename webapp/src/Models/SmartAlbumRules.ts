// SmartAlbumRules mirrors the backend SmartAlbumRules — the filter criteria
// a smart album's contents are derived from.
export interface SmartAlbumRules {
    startDate?: string;
    endDate?: string;
    mediaType?: string;
    camera?: string;
    hasGps?: boolean;
    orientation?: string;
    tags?: string[];
    favorite?: boolean;
    lowQuality?: boolean;
}
