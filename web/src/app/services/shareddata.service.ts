/** @format */

import { Injectable } from '@angular/core';
import { WSService } from '../api/ws/ws.service';
import { WSCommand, WSEvent } from '../api/ws/ws.static';
import {
  HelloEvent,
  JoinedEvent,
  PlayingEvent,
  EndEvent,
} from '../api/ws/ws.models';

@Injectable({
  providedIn: 'root',
})
export class SharedDataService {
  private _currentGuildID: string = null;

  constructor(ws: WSService) {
    ws.on(
      WSEvent.HELLO,
      (ev: HelloEvent) =>
        (this.currentGuildID = ev.voice_state ? ev.voice_state.guild_id : null)
    );

    ws.on(
      WSEvent.JOINED,
      (ev: JoinedEvent) => (this.currentGuildID = ev.guild_id)
    );

    ws.on(WSEvent.PLAYING, (ev: PlayingEvent) => {
      this.currentGuildID = ev.guild_id;
    });

    ws.on(WSEvent.END, (ev: EndEvent) => {
      this.currentGuildID = ev.guild_id;
    });
  }

  public set currentGuildID(val: string) {
    if (this.currentGuildID != val) this._currentGuildID = val;
  }

  public get currentGuildID(): string {
    return this._currentGuildID;
  }
}
