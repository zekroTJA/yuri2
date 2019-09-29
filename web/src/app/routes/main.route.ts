/** @format */

import { Component } from '@angular/core';
import { WSService } from '../api/ws/ws.service';
import { SoundListService } from '../services/soundlist.service';

interface SoundBtn {
  name: string;
}

@Component({
  selector: 'app-amin-route',
  templateUrl: './main.route.html',
  styleUrls: ['./main.route.sass'],
})
export class MainRouteComponent {
  constructor(private ws: WSService, private sounds: SoundListService) {}
}
