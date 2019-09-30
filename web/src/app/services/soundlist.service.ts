/** @format */

import { Injectable } from '@angular/core';
import { RestService } from '../api/rest/rest.service';
import { ActivatedRoute, Router } from '@angular/router';
import { Observable, of } from 'rxjs';
import { map } from 'rxjs/operators';

export interface SoundBtn {
  name: string;
  playing: boolean;
  favorite: boolean;
}

@Injectable({
  providedIn: 'root',
})
export class SoundListService {
  public sounds: SoundBtn[] = [];
  public sortBy: string;

  constructor(
    private rest: RestService,
    private route: ActivatedRoute,
    private router: Router
  ) {
    this.presetSortType.subscribe((sortBy) => {
      this.sortBy = sortBy;
      this.refreshSoundList();
    });
  }

  private get presetSortType(): Observable<string> {
    return this.route.queryParams.pipe(
      map((params) => {
        const p = params['sortBy'];
        if (p) return p;

        const ls = window.localStorage.getItem('sort_by');
        return ls || 'name';
      })
    );
  }

  public refreshSoundList() {
    this.rest.getSounds(this.sortBy, 0, 1000).subscribe((sounds) => {
      this.sounds = sounds.map<SoundBtn>(
        (s) =>
          ({
            name: s.name,
          } as SoundBtn)
      );
    });
    window.localStorage.setItem('sort_by', this.sortBy);
  }

  public setPlayingState(ident: string, playing: boolean) {
    const sound = this.sounds.find((s) => s.name == ident);
    if (sound) sound.playing = playing;
  }
}
