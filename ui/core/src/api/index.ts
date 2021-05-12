// Everything here represents services that are effectively remote data storage
export * from "./EthereumService/utils/getMetamaskProvider";

import ethereumService, { EthereumServiceContext } from "./EthereumService";
import ethbridgeService, { EthbridgeServiceContext } from "./EthbridgeService";
import sifService, { SifServiceContext } from "./SifService";
import clpService, { ClpServiceContext } from "./ClpService";
import eventBusService, { EventBusServiceContext } from "./EventBusService";
import createDispensationService, {
  IDispensationServiceContext,
} from "./DispensationService";

export type Api = ReturnType<typeof createApi>;

export type WithApi<T extends keyof Api = keyof Api> = {
  api: Pick<Api, T>;
};

export type ApiContext = EthereumServiceContext &
  SifServiceContext &
  ClpServiceContext &
  EthbridgeServiceContext &
  ClpServiceContext &
  EventBusServiceContext &
  IDispensationServiceContext; // add contexts from other APIs

export function createApi(context: ApiContext) {
  const EthereumService = ethereumService(context);
  const EthbridgeService = ethbridgeService(context);
  const SifService = sifService(context);
  const ClpService = clpService(context);
  const EventBusService = eventBusService(context);
  const DispensationService = createDispensationService(context);

  return {
    ClpService,
    EthereumService,
    SifService,
    EthbridgeService,
    EventBusService,
    DispensationService,
  };
}
