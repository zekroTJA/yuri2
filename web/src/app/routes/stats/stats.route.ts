/** @format */

import { Component } from '@angular/core';
import { RestService } from 'src/app/api/rest/rest.service';
import { StatsEntry } from 'src/app/api/rest/rest.models';
import { SharedDataService } from 'src/app/services/shareddata.service';
import { ToastService } from 'src/app/components/toast/toast.service';
import { WSService } from 'src/app/api/ws/ws.service';
import dateFormat from 'dateformat';

@Component({
  selector: 'app-stats-route',
  templateUrl: './stats.route.html',
  styleUrls: ['./stats.route.sass'],
})
export class StatsRouteComponent {
  public stats: StatsEntry[];
  public dateFormat = dateFormat;

  constructor(
    private rest: RestService,
    private sharedData: SharedDataService,
    private ws: WSService,
    private toasts: ToastService
  ) {
    this.fetchStats();
  }

  public fetchStats() {
    if (!this.ws.isInitialized) {
      setTimeout(this.fetchStats.bind(this), 100);
      return;
    }

    if (!this.sharedData.currentGuildID) {
      this.toasts.push(
        "You need to be in a Voice Channel with Yuri that the Guild's stats can be fetched!",
        'Error',
        'error',
        10000,
        true
      );
      return;
    }

    this.rest
      .getPlayStats(this.sharedData.currentGuildID, 100)
      .subscribe((stats) => (this.stats = stats));
  }

  public getSourceClassName(source: string) {
    return `source-${source.toLowerCase()}`;
  }
}
