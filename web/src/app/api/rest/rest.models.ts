/** @format */

export interface ListReponse<T> {
  n: number;
  results: T[];
}

export interface Sound {
  name: string;
  last_modified: Date;
}

export interface LogEntry {
  time: Date;
  user_id: string;
  user_tag: string;
  guild_id: string;
  source: string;
  sound: string;
}

export interface StatsEntry {
  sound: string;
  count: number;
}

export interface FastTrigger {
  ident: string;
  random: boolean;
}

export interface Guild {
  name: string;
  id: string;
}

export interface VoiceConnection {
  guild: Guild;
  vc_id: string;
}

export interface SystemDetails {
  arch: string;
  os: string;
  go_version: string;
  cpu_used_cores: number;
  go_routines: number;
  heap_use_b: number;
  stack_use_b: number;
  uptime_seconds: number;
}

export interface SystemStats {
  guilds: Guild[];
  voice_connections: VoiceConnection[];
  system: SystemDetails;
}

export interface SoundStats {
  sounds_len: number;
  log_len: number;
  size_b: number;
}

export interface YouTubeEmbed {
  author_name: string;
  title: string;
  thumbnail_url: string;
  error: string;
}

export interface BuildInfo {
  build_version: string;
  go_version: string;
  commit_hash: string;
  go_os: string;
  go_arch: string;
}
