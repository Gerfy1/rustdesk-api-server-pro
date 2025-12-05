import { request } from '@/service/request';

// ========== CATEGORIES ==========

/** Get all KB categories */
export function fetchKBCategories() {
  return request<Api.DocHelp.CategoryList>({
    url: '/dochelp/categories',
    method: 'get'
  });
}

/** Create KB category */
export function createKBCategory(data: Api.DocHelp.CategoryEdit) {
  return request<Api.Common.CommonResult>({
    url: '/dochelp/categories',
    method: 'post',
    data
  });
}

/** Update KB category */
export function updateKBCategory(id: number, data: Api.DocHelp.CategoryEdit) {
  return request<Api.Common.CommonResult>({
    url: `/dochelp/categories/${id}`,
    method: 'put',
    data
  });
}

/** Delete KB category */
export function deleteKBCategory(id: number) {
  return request<Api.Common.CommonResult>({
    url: `/dochelp/categories/${id}`,
    method: 'delete'
  });
}

// ========== ARTICLES ==========

/** Get KB articles with pagination */
export function fetchKBArticles(params?: Api.DocHelp.ArticleParams) {
  return request<Api.DocHelp.ArticleList>({
    url: '/dochelp/articles',
    method: 'get',
    params
  });
}

/** Get single KB article */
export function fetchKBArticle(id: number) {
  return request<Api.DocHelp.ArticleDetail>({
    url: `/dochelp/articles/${id}`,
    method: 'get'
  });
}

/** Create KB article */
export function createKBArticle(data: Api.DocHelp.ArticleEdit) {
  return request<Api.Common.CommonResult>({
    url: '/dochelp/articles',
    method: 'post',
    data
  });
}

/** Update KB article */
export function updateKBArticle(id: number, data: Api.DocHelp.ArticleEdit) {
  return request<Api.Common.CommonResult>({
    url: `/dochelp/articles/${id}`,
    method: 'put',
    data
  });
}

/** Delete KB article */
export function deleteKBArticle(id: number) {
  return request<Api.Common.CommonResult>({
    url: `/dochelp/articles/${id}`,
    method: 'delete'
  });
}

/** Increment article view count */
export function incrementArticleView(id: number) {
  return request<Api.Common.CommonResult>({
    url: `/dochelp/articles/${id}/view`,
    method: 'post'
  });
}

// ========== TICKETS ==========

/** Get tickets with pagination */
export function fetchTickets(params?: Api.DocHelp.TicketParams) {
  return request<Api.DocHelp.TicketList>({
    url: '/dochelp/tickets',
    method: 'get',
    params
  });
}

/** Get single ticket */
export function fetchTicket(id: number) {
  return request<Api.DocHelp.TicketDetail>({
    url: `/dochelp/tickets/${id}`,
    method: 'get'
  });
}

/** Create ticket */
export function createTicket(data: Api.DocHelp.TicketCreate) {
  return request<Api.Common.CommonResult>({
    url: '/dochelp/tickets',
    method: 'post',
    data
  });
}

/** Update ticket */
export function updateTicket(id: number, data: Api.DocHelp.TicketUpdate) {
  return request<Api.Common.CommonResult>({
    url: `/dochelp/tickets/${id}`,
    method: 'put',
    data
  });
}

/** Resolve ticket */
export function resolveTicket(id: number) {
  return request<Api.Common.CommonResult>({
    url: `/dochelp/tickets/${id}/resolve`,
    method: 'post'
  });
}

/** Delete ticket */
export function deleteTicket(id: number) {
  return request<Api.Common.CommonResult>({
    url: `/dochelp/tickets/${id}`,
    method: 'delete'
  });
}

/** Convert ticket to KB article */
export function convertTicketToArticle(id: number, data: Api.DocHelp.TicketConvert) {
  return request<Api.Common.CommonResult>({
    url: `/dochelp/tickets/${id}/convert`,
    method: 'post',
    data
  });
}

// ========== COMMENTS ==========

/** Get ticket comments */
export function fetchTicketComments(ticketId: number) {
  return request<Api.DocHelp.CommentList>({
    url: `/dochelp/tickets/${ticketId}/comments`,
    method: 'get'
  });
}

/** Add comment to ticket */
export function addTicketComment(ticketId: number, data: Api.DocHelp.CommentCreate) {
  return request<Api.Common.CommonResult>({
    url: `/dochelp/tickets/${ticketId}/comments`,
    method: 'post',
    data
  });
}

/** Delete comment */
export function deleteTicketComment(commentId: number) {
  return request<Api.Common.CommonResult>({
    url: `/dochelp/comments/${commentId}`,
    method: 'delete'
  });
}

// ========== STATISTICS ==========

/** Get DocHelp statistics */
export function fetchDocHelpStats() {
  return request<Api.DocHelp.Stats>({
    url: '/dochelp/stats',
    method: 'get'
  });
}

// ========== FILE UPLOAD ==========

/** Upload file for DocHelp */
export function uploadDocHelpFile(file: File) {
  const formData = new FormData();
  formData.append('file', file);
  
  return request<{ filename: string; url: string; size: number }>({
    url: '/dochelp/upload',
    method: 'post',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  });
}
