/** @format */

import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { SideBarComponent } from './components/sidebar/sidebar.component';
import { SliderComponent } from './components/slider/slider.component';
import { FormsModule } from '@angular/forms';
import { ToastComponent } from './components/toast/toast.component';
import { MainRouteComponent } from './routes/main/main.route';
import { SoundBtnComponent } from './components/soundbtn/soundbtn.component';
import { HttpClientModule } from '@angular/common/http';
import { SearchBarComponent } from './components/searchbar/searchbar.component';
import { ContextMenuComponent } from './components/contextmenu/contextmenu.component';
import { LogsRouteComponent } from './routes/logs/logs.route';
import { StatsRouteComponent } from './routes/stats/stats.route';
import { AdminRouteComponent } from './routes/admin/admin.route';
import { YTInputComponent } from './components/ytimport/ytimport.component';
import { InfoRouteComponent } from './routes/info/info.route';

@NgModule({
  declarations: [
    AppComponent,

    MainRouteComponent,
    LogsRouteComponent,
    StatsRouteComponent,
    AdminRouteComponent,
    InfoRouteComponent,

    SideBarComponent,
    SliderComponent,
    ToastComponent,
    SoundBtnComponent,
    SearchBarComponent,
    ContextMenuComponent,
    YTInputComponent,
  ],
  imports: [BrowserModule, AppRoutingModule, FormsModule, HttpClientModule],
  providers: [],
  bootstrap: [AppComponent],
})
export class AppModule {}
