/** @format */

import { Injectable } from '@angular/core';
import { RestService } from '../api/rest/rest.service';

interface SoundBtn {
  name: string;
}

@Injectable({
  providedIn: 'root',
})
export class SoundListService {
  public sounds: SoundBtn[] = [];

  constructor(private rest: RestService) {
    this.refreshSoundList();
  }

  public refreshSoundList(sortBy: string = 'name') {
    this.rest.getSounds(sortBy, 0, 1000).subscribe((sounds) => {
      this.sounds = sounds.map<SoundBtn>(
        (s) =>
          ({
            name: s.name,
          } as SoundBtn)
      );
    });
  }
}
