import pm2 from "pm2"

export type PM2Process = string | number;

/**
 * Utility Functions that extend PM2
 */

export async function spawn(command: string, args: string[]) {
  await connect(false)
  await start({
    script: command,
    args: args,
    autorestart: false,
    force: true
  } as pm2.StartOptions);
}

/**
 * Promisified Functions Below
 */

export const connect = function (no_daemon_mode: boolean) {
  return new Promise<void>(function (resolve, reject) {
    pm2.connect(no_daemon_mode, function (err) {
      if (err) {
        return reject(err)
      }
      return resolve()
    });
  })
}

export const disconnect = pm2.disconnect

export const start = function (process: pm2.StartOptions) {
  return new Promise<void>(function (resolve, reject) {
    pm2.start(process, function (err) {
      if (err) {
        return reject(err);
      }
      return resolve();
    });
  })
}

export const stop = function (process: PM2Process) {
  return new Promise<void>(function (resolve, reject) {
    pm2.stop(process, function (err) {
      if (err) {
        return reject(err);
      }
      return resolve();
    })
  })
}

export const restart = function (process: PM2Process) {
  return new Promise<void>(function (resolve, reject) {
    pm2.restart(process, function (err) {
      if (err) {
        return reject(err);
      }
      return resolve();
    })
  })
}

export const reload = function (process: PM2Process) {
  return new Promise<void>(function (resolve, reject) {
    pm2.reload(process, function (err) {
      if (err) {
        return reject(err)
      }
      return resolve();
    })
  })
}

export const del = function (process: PM2Process) {
  return new Promise<void>(function (resolve, reject) {
    pm2.delete(process, function (err) {
      if (err) {
        return reject(err);
      }
      return resolve();
    })
  })
}

export const killDaemon = function () {
  return new Promise<void>(function (resolve, reject) {
    pm2.killDaemon(function (err) {
      if (err) {
        return reject(err);
      }
      return resolve();
    })
  })
}

export const describe = function (process: PM2Process) {
  return new Promise<void>(function (resolve, reject) {
    pm2.describe(process, function (err) {
      if (err) {
        return reject(err);
      }
      return resolve();
    })
  })
}

export const list = function () {
  return new Promise(function (resolve, reject) {
    pm2.list(function (err, list) {
      if (err) {
        return reject(err);
      }
      return resolve(list)
    })
  })
}

export interface PM2DataPacket {
  id: number;
  type: string;
  topic: boolean;
  data: unknown;
}

export const sendDataToProcessId = function (proc_id: number, packet: PM2DataPacket) {
  return new Promise(function (reslove, reject) {
    pm2.sendDataToProcessId(proc_id, packet, function (err, res) {
      if (err) {
        return reject(err);
      }
      return reslove(res);
    })
  })
}

export const sendSignalToProcessName = function (signal: string | number, process: string | number) {
  return new Promise<void>(function (resolve, reject) {
    pm2.sendSignalToProcessName(signal, process, function (err, res) {
      if (err) {
        return reject(err);
      }
      return resolve(res)
    })
  })
}
