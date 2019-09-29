/** @format */

import { Injectable } from '@angular/core';
import { environment } from 'src/environments/environment';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Observable, of } from 'rxjs';
import { catchError, map } from 'rxjs/operators';
import {
  ListReponse,
  Sound,
  LogEntry,
  StatsEntry,
  FastTrigger,
  SystemStats,
  SoundStats,
} from './rest.models';
import { ToastService } from 'src/app/components/toast/toast.service';

/** @format */

@Injectable({
  providedIn: 'root',
})
export class RestService {
  private rootURL = '';

  private readonly errorCatcher = (err) => {
    console.error(err);
    this.toasts.push(err.message, 'REST Request Error', 'error', 10000);
    return of(null);
  };

  private readonly defopts = (obj?: object) => {
    const defopts = {
      withCredentials: true,
    };

    if (obj) {
      Object.keys(obj).forEach((k) => {
        defopts[k] = obj[k];
      });
    }

    return defopts;
  };

  // ----------------
  // RESOURCES

  private readonly rcRoot = (sub: string = null) =>
    this.rootURL + '/api' + (sub ? `/${sub}` : '');

  private readonly rcPlayLogs = (guild: string) =>
    this.rcRoot('logs') + '/' + guild;

  private readonly rcPlayStats = (guild: string) =>
    this.rcRoot('stats') + '/' + guild;

  private readonly rcFavorites = (sub: string = null) =>
    this.rcRoot('favorite') + (sub ? `/${sub}` : '');

  private readonly rcSettings = (sub: string = null) =>
    this.rcRoot('settings') + (sub ? `/${sub}` : '');

  private readonly rcAdmin = (sub: string = null) =>
    this.rcRoot('admin') + (sub ? `/${sub}` : '');

  // ----------------

  constructor(private http: HttpClient, private toasts: ToastService) {
    this.rootURL = environment.production ? '' : 'http://localhost:8080';
  }

  public getSounds(
    sort: string,
    from: number = 0,
    limit: number = 0
  ): Observable<Sound[]> {
    const opts = this.defopts({
      params: new HttpParams()
        .set('sort', sort)
        .set('from', from.toString())
        .set('limit', limit.toString()),
    });
    return this.http
      .get<ListReponse<Sound>>(this.rcRoot('localsounds'), opts)
      .pipe(
        map((lr) => lr.results),
        catchError(this.errorCatcher)
      );
  }

  public getPlayLog(
    guild: string,
    from: number,
    limit: number
  ): Observable<LogEntry[]> {
    const opts = this.defopts({
      params: new HttpParams()
        .set('from', from.toString())
        .set('limit', limit.toString()),
    });
    return this.http
      .get<ListReponse<LogEntry>>(this.rcPlayLogs(guild), opts)
      .pipe(
        map((lr) => lr.results),
        catchError(this.errorCatcher)
      );
  }

  public getPlayStats(guild: string, limit: number): Observable<StatsEntry[]> {
    const opts = this.defopts({
      params: new HttpParams().set('limit', limit.toString()),
    });
    return this.http
      .get<ListReponse<StatsEntry>>(this.rcPlayStats(guild), opts)
      .pipe(
        map((lr) => lr.results),
        catchError(this.errorCatcher)
      );
  }

  public getFavorites(): Observable<string[]> {
    return this.http
      .get<ListReponse<string>>(this.rcFavorites(), this.defopts())
      .pipe(
        map((lr) => lr.results),
        catchError(this.errorCatcher)
      );
  }

  public setFavorite(sound: string): Observable<any> {
    return this.http
      .post(this.rcFavorites(sound), this.defopts())
      .pipe(catchError(this.errorCatcher));
  }

  public deleteFavorite(sound: string): Observable<any> {
    return this.http
      .delete(this.rcFavorites(sound), this.defopts())
      .pipe(catchError(this.errorCatcher));
  }

  public getFastTrigger(): Observable<FastTrigger> {
    return this.http
      .get<FastTrigger>(this.rcSettings('fasttrigger'), this.defopts())
      .pipe(catchError(this.errorCatcher));
  }

  public setFastTrigger(ident: string, random: boolean): Observable<any> {
    return this.http
      .post(this.rcSettings('fasttrigger'), { ident, random }, this.defopts())
      .pipe(catchError(this.errorCatcher));
  }

  public getSystemStats(): Observable<SystemStats> {
    return this.http
      .get<SystemStats>(this.rcAdmin('stats'), this.defopts())
      .pipe(catchError(this.errorCatcher));
  }

  public getSoundStats(): Observable<SoundStats> {
    return this.http
      .get<SoundStats>(this.rcAdmin('soundstats'), this.defopts())
      .pipe(catchError(this.errorCatcher));
  }

  public postRestart(): Observable<any> {
    return this.http
      .post(this.rcAdmin('restart'), this.defopts())
      .pipe(catchError(this.errorCatcher));
  }

  public postRefetch(): Observable<any> {
    return this.http
      .post(this.rcAdmin('refetch'), this.defopts())
      .pipe(catchError(this.errorCatcher));
  }
}
