/** @format */

import { Component, OnInit, OnDestroy } from '@angular/core';
import { WSService } from '../api/ws/ws.service';
import { SoundListService, SoundBtn } from '../services/soundlist.service';
import { WSCommand, SourceType, WSEvent } from '../api/ws/ws.static';
import { PlayingEvent, EndEvent } from '../api/ws/ws.models';

@Component({
  selector: 'app-main-route',
  templateUrl: './main.route.html',
  styleUrls: ['./main.route.sass'],
})
export class MainRouteComponent implements OnInit, OnDestroy {
  public search: boolean;
  public displayedSounds: SoundBtn[];

  constructor(private ws: WSService, private sounds: SoundListService) {
    this.ws.on(WSEvent.PLAYING, (ev: PlayingEvent) => {
      sounds.setPlayingState(ev.ident, true);
    });

    this.ws.on(WSEvent.END, (ev: EndEvent) => {
      sounds.setPlayingState(ev.ident, false);
    });
  }

  public playSound(sound: SoundBtn) {
    this.ws.sendMessage(WSCommand.PLAY, {
      ident: sound.name,
      source: SourceType.LOCAL,
    });
  }

  private onKeyDown(event: any) {
    if (event.keyCode == 114 || (event.ctrlKey && event.keyCode == 70)) {
      this.search = true;
      event.preventDefault();
    } else if (event.keyCode == 27) {
      this.onSearchClose();
      event.preventDefault();
    }
  }

  public onSearchInput(val: string) {
    this.displayedSounds = this.sounds.sounds.filter((s) =>
      s.name.toLowerCase().includes(val.toLowerCase())
    );
  }

  public onSearchClose() {
    this.search = false;
    this.displayedSounds = null;
  }

  public ngOnInit() {
    window.addEventListener('keydown', this.onKeyDown.bind(this));
  }

  public ngOnDestroy() {
    window.removeEventListener('keydown', this.onKeyDown.bind(this));
  }
}
