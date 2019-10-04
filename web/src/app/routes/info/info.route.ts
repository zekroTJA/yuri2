/** @format */

import { Component } from '@angular/core';
import { RestService } from 'src/app/api/rest/rest.service';
import { BuildInfo } from 'src/app/api/rest/rest.models';

@Component({
  selector: 'app-info-route',
  templateUrl: './info.route.html',
  styleUrls: ['./info.route.sass'],
})
export class InfoRouteComponent {
  public info: BuildInfo;

  constructor(private rest: RestService) {
    rest.getInfo().subscribe((info) => (this.info = info));
  }
}
