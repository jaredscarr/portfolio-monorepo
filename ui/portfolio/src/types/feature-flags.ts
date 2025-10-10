export interface FeatureFlag {
  key: string;
  enabled: boolean;
}

export interface FlagsResponse {
  [key: string]: boolean;
}

export interface UpdateFlagRequest {
  enabled: boolean;
}

export interface FlagUpdateResponse {
  key: string;
  enabled: boolean;
}

export interface ReloadResponse {
  status: string;
}

export interface ErrorResponse {
  error: string;
}
