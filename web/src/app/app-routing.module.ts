/** @format */

import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { MainRouteComponent } from './routes/main.route';

const routes: Routes = [
  {
    path: 'sounds',
    component: MainRouteComponent,
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
