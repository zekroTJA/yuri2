/** @format */

import { Component } from '@angular/core';
import { RestService } from 'src/app/api/rest/rest.service';
import { SystemStats } from 'src/app/api/rest/rest.models';
import { ToastService } from 'src/app/components/toast/toast.service';
import dateFormat from 'dateformat';
import { toDDHHMMSS } from 'src/util/util.time';
import { byteCountFormatter } from 'src/util/util.format';

@Component({
  selector: 'app-admin-route',
  templateUrl: './admin.route.html',
  styleUrls: ['./admin.route.sass'],
})
export class AdminRouteComponent {
  public stats: SystemStats;
  public uptime: number;
  public dateFormat = dateFormat;

  public toDDHHMMSS = toDDHHMMSS;
  public byteCountFormatter = byteCountFormatter;

  constructor(private rest: RestService, private toasts: ToastService) {
    this.fetchData().then(() => setInterval(() => this.uptime++, 1000));
  }

  public fetchData(): Promise<any> {
    return new Promise((resolve, reject) => {
      this.rest.getSystemStats().subscribe((stats: SystemStats) => {
        this.stats = stats;
        this.uptime = stats.system.uptime_seconds;
        resolve();
      });
    });
  }

  public refetch() {
    this.rest
      .postRefetch()
      .toPromise()
      .then(() => {
        this.toasts.push(
          'Refetched sound list',
          'Refetch',
          'success',
          6000,
          true
        );
      });
  }

  public restart() {
    this.rest
      .postRestart()
      .toPromise()
      .then(() => {
        this.toasts.push(
          'This site will reload in 5 seconds which is the estimated time until the backend should be online again.',
          'Restarting...',
          'success',
          6000,
          true
        );

        setTimeout(() => window.location.reload(), 5000);
      });
  }
}
