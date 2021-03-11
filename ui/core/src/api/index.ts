// Everything here represents services that are effectively remote data storage
export * from "./EthereumService/utils/getMetamaskProvider";

import ethereumService, { EthereumServiceContext } from "./EthereumService";
import ethbridgeService, { EthbridgeServiceContext } from "./EthbridgeService";
import sifService, { SifServiceContext } from "./SifService";
import clpService, { ClpServiceContext } from "./ClpService";
import notificationService, { EventsServiceContext } from "./EventsService";

export type Api = ReturnType<typeof createApi>;

export type WithApi<T extends keyof Api = keyof Api> = {
  api: Pick<Api, T>;
};

export type ApiContext = EthereumServiceContext &
  SifServiceContext &
  ClpServiceContext &
  EthbridgeServiceContext &
  ClpServiceContext &
  EventsServiceContext; // add contexts from other APIs

export function createApi(context: ApiContext) {
  const EthereumService = ethereumService(context);
  const EthbridgeService = ethbridgeService(context);
  const SifService = sifService(context);
  const ClpService = clpService(context);
  const EventsService = notificationService(context);
  return {
    ClpService,
    EthereumService,
    SifService,
    EthbridgeService,
    EventsService,
  };
}
