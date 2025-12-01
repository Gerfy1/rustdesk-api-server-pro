import { request } from '../request';

export function fetchDevicesList(params: any) {
  return request<Api.Devices.DevicesList>({ url: '/devices/list', params });
}

export function fetchDevicesOnlineList() {
  return request<Api.Devices.Device[]>({ url: '/devices/online' });
}
