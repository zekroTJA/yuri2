/** @format */

import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { MainRouteComponent } from './routes/main/main.route';
import { LogsRouteComponent } from './routes/logs/logs.route';
import { StatsRouteComponent } from './routes/stats/stats.route';
import { AdminRouteComponent } from './routes/admin/admin.route';

const routes: Routes = [
  {
    path: 'sounds',
    component: MainRouteComponent,
  },
  {
    path: 'logs',
    component: LogsRouteComponent,
  },
  {
    path: 'stats',
    component: StatsRouteComponent,
  },
  {
    path: 'admin',
    component: AdminRouteComponent,
  },
  {
    path: '**',
    redirectTo: '/sounds',
    pathMatch: 'full',
  },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule],
})
export class AppRoutingModule {}
