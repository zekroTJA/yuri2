/** @format */

import { Component } from '@angular/core';
import { WSService } from 'src/app/api/ws/ws.service';
import { WSEvent, WSCommand } from 'src/app/api/ws/ws.static';
import {
  HelloEvent,
  JoinedEvent,
  LeftEvent,
  VolumeChangedEvent,
} from 'src/app/api/ws/ws.models';
import { toNumber } from 'src/util/util.converters';
import { SoundListService } from 'src/app/services/soundlist.service';

@Component({
  selector: 'app-sidebar',
  templateUrl: './sidebar.component.html',
  styleUrls: ['./sidebar.component.sass'],
})
export class SideBarComponent {
  public inChannel: boolean;
  public isAdmin: boolean;
  public volume: number;

  constructor(private ws: WSService, public sounds: SoundListService) {
    ws.on(WSEvent.HELLO, (ev: HelloEvent) => {
      this.isAdmin = ev.admin;
      this.inChannel = ev.connected;
      this.volume = ev.vol;
    });

    ws.on(WSEvent.JOINED, (ev: JoinedEvent) => {
      this.inChannel = true;
      this.volume = ev.vol;
    });

    ws.on(WSEvent.LEFT, (ev: LeftEvent) => {
      this.inChannel = false;
    });

    ws.on(WSEvent.VOLUME_CHANGED, (ev: VolumeChangedEvent) => {
      this.volume = ev.vol;
    });
  }

  public onSortBy() {
    this.sounds.sortBy = this.sounds.sortBy == 'name' ? 'date' : 'name';
    this.sounds.refreshSoundList();
  }

  public onStop() {
    this.ws.sendMessage(WSCommand.STOP);
  }

  public onJoin() {
    const cmd = this.inChannel ? WSCommand.LEAVE : WSCommand.JOIN;
    this.ws.sendMessage(cmd);
  }

  public onVolume() {
    this.ws.sendMessage(WSCommand.VOLUME, this.volume);
  }
}
