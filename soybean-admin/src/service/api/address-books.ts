import { request } from '../request';

/** Fetch address books list */
export function fetchAddressBooksList(params: any) {
  return request<Api.AddressBooks.AddressBooksList>({ url: '/address-books/list', params });
}

/** Fetch single address book */
export function fetchAddressBook(id: number) {
  return request<Api.AddressBooks.AddressBookDetail>({ url: `/address-books/${id}` });
}

/** Create address book */
export function createAddressBook(data: any) {
  return request({ url: '/address-books', method: 'post', data });
}

/** Update address book */
export function updateAddressBook(id: number, data: any) {
  return request({ url: `/address-books/${id}`, method: 'put', data });
}

/** Delete address book */
export function deleteAddressBook(id: number) {
  return request({ url: `/address-books/${id}`, method: 'delete' });
}

/** Fetch peers in address book */
export function fetchAddressBookPeers(id: number) {
  return request<Api.AddressBooks.PeersList>({ url: `/address-books/${id}/peers` });
}

/** Import devices as peers */
export function importDevicesAsPeers(id: number) {
  return request({ url: `/address-books/${id}/import-devices`, method: 'post' });
}
