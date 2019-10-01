/** @format */

import { SourceType } from './ws.static';

/** @format */

export interface Message {
  name: string;
  data: any;
}

export interface VoiceState {
  guild_id: string;
  channel_id: string;
}

export interface Track {
  ident: string;
  source: SourceType;
  guild_id: string;
  channel_id: string;
  user_id: string;
  user_tag: string;
}

export interface HelloEvent {
  admin: boolean;
  connected: boolean;
  vol: number;
  voice_state: VoiceState;
}

export interface WSErrorEvent {
  code: number;
  type: string;
  message: string;
  data: any;
}

export interface PlayingEvent {
  ident: string;
  source: SourceType;
  guild_id: string;
  user_id: string;
  user_tag: string;
  vol: number;
}

export interface EndEvent {
  ident: string;
  source: SourceType;
  guild_id: string;
  user_id: string;
  user_tag: string;
}

export interface PlayErrorEvent {
  reason: string;
  track: Track;
}

export interface StuckEvent {
  threshold: number;
  track: Track;
}

export interface VolumeChangedEvent {
  guild_id: string;
  vol: number;
}

export interface JoinedEvent {
  guild_id: string;
  channel_id: string;
  vol: number;
}

export interface LeftEvent {
  guild_id: string;
  channel_id: string;
}
