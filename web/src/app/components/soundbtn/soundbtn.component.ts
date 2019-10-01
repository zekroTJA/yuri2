/** @format */

import { Component, Input } from '@angular/core';

@Component({
  selector: 'app-soundbtn',
  templateUrl: './soundbtn.component.html',
  styleUrls: ['./soundbtn.component.sass'],
})
export class SoundBtnComponent {
  @Input() public name: string;
  @Input() public playing: boolean;
  @Input() public favorite: boolean;
}
