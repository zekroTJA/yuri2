/** @format */

import { HelloEvent } from './ws.models';

export enum SourceType {
  LOCAL,
  YOUTUBE,
  HTTPURI,
}

export enum ErrorType {
  BAD_COMMAND_ARGS,
  UNAUTHORIZED,
  FORBIDDEN,
  INTERNAL,
  BAD_COMMAND,
  RATE_LIMIT_EXCEED,
}

export enum WSCommand {
  INIT = 'INIT',
  JOIN = 'JOIN',
  LEAVE = 'LEAVE',
  PLAY = 'PLAY',
  RANDOM = 'RANDOM',
  STOP = 'STOP',
  VOLUME = 'VOLUME',
}

export enum WSEvent {
  HELLO = 'HELLO',
  ERROR = 'ERROR',
  PLAYING = 'PLAYING',
  END = 'END',
  PLAY_ERROR = 'PLAY_ERROR',
  STUCK = 'STUCK',
  VOLUME_CHANGED = 'VOLUME_CHANGED',
  JOINED = 'JOINED',
  LEFT = 'LEFT',
}
