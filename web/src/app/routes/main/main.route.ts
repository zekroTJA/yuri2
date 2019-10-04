/** @format */

import { Component, OnInit, OnDestroy } from '@angular/core';
import { WSService } from '../../api/ws/ws.service';
import { SoundListService, SoundBtn } from '../../services/soundlist.service';
import { WSCommand, SourceType, WSEvent } from '../../api/ws/ws.static';
import { PlayingEvent, EndEvent } from '../../api/ws/ws.models';
import { ContextMenuItem } from '../../components/contextmenu/contextmenu.component';
import { RestService } from '../../api/rest/rest.service';
import { rand } from 'src/util/util.random';

@Component({
  selector: 'app-main-route',
  templateUrl: './main.route.html',
  styleUrls: ['./main.route.sass'],
})
export class MainRouteComponent implements OnInit, OnDestroy {
  public search: boolean;
  public displayedSounds: SoundBtn[];

  public randSkeletonWidths: number[] = Array(100)
    .fill(1)
    .map((n) => rand(50, 100));

  public contextMenu = {
    x: 0,
    y: 0,
    visible: false,
  };

  public contextMenuItems: ContextMenuItem[] = [
    {
      el: 'Favorite',
      action: () => {
        this.contextMenu.visible = false;
      },
    },
  ];

  constructor(
    private rest: RestService,
    private ws: WSService,
    public sounds: SoundListService
  ) {
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
    if (event.ctrlKey && event.keyCode == 70) {
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

  public onSoundContextMenu(ev: any, sound: SoundBtn) {
    this.contextMenuItems[0].el = sound.favorite ? 'Unfavorize' : 'Favorize';
    this.contextMenuItems[0].action = sound.favorite
      ? this.unfavorize.bind(this, sound)
      : this.favorize.bind(this, sound);

    this.contextMenu.visible = true;
    this.contextMenu.x = ev.clientX;
    this.contextMenu.y = ev.clientY;

    ev.preventDefault();
  }

  private favorize(sound: SoundBtn) {
    this.rest
      .setFavorite(sound.name)
      .toPromise()
      .then(() => (sound.favorite = true));
  }

  private unfavorize(sound: SoundBtn) {
    this.rest
      .deleteFavorite(sound.name)
      .toPromise()
      .then(() => (sound.favorite = false));
  }

  public onWindowClick(ev: any) {
    if (ev.target && ev.target.id != 'context-menu') {
      this.contextMenu.visible = false;
    }
  }

  public ngOnInit() {
    window.addEventListener('keydown', this.onKeyDown.bind(this));
    window.addEventListener('click', this.onWindowClick.bind(this));
  }

  public ngOnDestroy() {
    window.removeEventListener('keydown', this.onKeyDown.bind(this));
    window.removeEventListener('click', this.onWindowClick.bind(this));
  }
}
