export interface OutboxEvent {
  id: string;
  type: string;
  source: string;
  data: Record<string, unknown>;
  metadata?: Record<string, unknown>;
  status: "pending" | "published" | "failed" | "retrying";
  error_message?: string;
  retry_count: number;
  created_at: string;
  published_at?: string;
}

export interface EventsResponse {
  events: OutboxEvent[];
  total: number;
  page: number;
  limit: number;
}

export interface CreateEventRequest {
  type: string;
  source: string;
  data: Record<string, unknown>;
  metadata?: Record<string, unknown>;
}

export interface OutboxStats {
  total_events: number;
  pending_events: number;
  published_events: number;
  failed_events: number;
  retry_count: number;
}

export interface PublishRequest {
  batch_size?: number;
  event_ids?: string[];
}

export interface PublishResponse {
  published: number;
  failed: number;
  errors?: string[];
}
