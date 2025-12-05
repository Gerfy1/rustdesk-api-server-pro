/**
 * Namespace Api
 *
 * All backend api type
 */
declare namespace Api {
  namespace Common {
    /** common params of paginating */
    interface PaginatingCommonParams {
      /** current page number */
      current: number;
      /** page size */
      size: number;
      /** total count */
      total: number;
    }

    /** common params of paginating query list data */
    interface PaginatingQueryRecord<T = any> extends PaginatingCommonParams {
      records: T[];
    }

    /** common search params of table */
    type CommonSearchParams = Pick<Common.PaginatingCommonParams, 'current' | 'size'>;

    /** common record */
    type CommonRecord<T = any> = {
      id?: number;
      created_at?: string;
    } & T;

    /** common result for operations */
    type CommonResult = {
      data: any;
      error: string | null;
    };
  }

  namespace Form {
    interface LoginForm {
      username: string;
      password: string;
      code: string;
      captchaId: string;
    }
  }

  /**
   * namespace Auth
   *
   * backend api module: "auth"
   */
  namespace Auth {
    interface LoginToken {
      token: string;
    }

    interface Captcha {
      id: string;
      img: string;
    }

    interface UserInfo {
      userId: string;
      userName: string;
      role?: number; // User role: 1=User, 2=Support, 3=Support N2, 4=Super Admin
      roles: string[];
      buttons: string[];
    }
  }

  namespace Home {
    interface Stat {
      userCount: number;
      deviceCount: number;
      onlineCount: number;
      visitsCount: number;
    }

    interface LineChart {
      xAxis: string[];
      users: number[];
      peer: number[];
    }

    interface PieChart {
      name: string;
      value: number;
    }
  }

  namespace UserManagement {
    type User = Common.CommonRecord<{
      username: string;
      password: string;
      name: string;
      email: string;
      licensed_devices: number;
      login_verify: string;
      tfa_secret: string;
      note: string;
      status: number;
      is_admin: boolean;
      admin_status: number;
      tfa_code: string;
      role: number; // 1=User, 2=Support, 3=Support N2, 4=Super Admin
    }>;
    type UserList = Common.PaginatingQueryRecord<User>;

    type UserSearchParams = CommonType.RecordNullable<
      Pick<Api.UserManagement.User, 'username' | 'name' | 'email' | 'status' | 'admin_status' | 'created_at'> &
        Api.Common.CommonSearchParams
    >;

    type Session = Common.CommonRecord<{
      username: string;
      rustdesk_id: string;
      expired: string;
    }>;
    type SessionList = Common.PaginatingQueryRecord<Session>;

    type SessionSearchParams = CommonType.RecordNullable<
      Pick<Api.UserManagement.Session, 'username' | 'created_at'> & Api.Common.CommonSearchParams
    >;
  }

  namespace Devices {
    type Device = Common.CommonRecord<{
      rustdesk_id: string;
      hostname: string;
      username: string;
      uuid: string;
      version: string;
      os: string;
      memory: string;
      is_online: boolean;
      last_seen_at: string;
      ip_address: string;
      conns: number;
      created_at: string;
      total_accesses?: number;
      last_connection_at?: string;
    }>;
    type DevicesList = Common.PaginatingQueryRecord<Device>;
    type DeviceSearchParams = CommonType.RecordNullable<
      Pick<
        Api.Devices.Device,
        'username' | 'hostname' | 'rustdesk_id'
      > &
        Api.Common.CommonSearchParams & {
          status?: string;
        }
    >;
  }

  namespace Audit {
    type AuditLog = Common.CommonRecord<{
      username: string;
      conn_id: string;
      rustdesk_id: string;
      ip: string;
      session_id: string;
      uuid: string;
      type: number;
      closed_at?: string;
      duration?: string; // Calculated field in frontend
    }>;
    type AuditLogList = Common.PaginatingQueryRecord<AuditLog>;
    type AuditLogSearchParams = CommonType.RecordNullable<
      Pick<
        Api.Audit.AuditLog,
        'username' | 'conn_id' | 'rustdesk_id' | 'ip' | 'session_id' | 'uuid' | 'type' | 'created_at' | 'closed_at'
      > &
        Api.Common.CommonSearchParams
    >;

    type AuditFileTransferLog = Common.CommonRecord<{
      rustdesk_id: string;
      peer_id: string;
      path: string;
      uuid: string;
      type: number;
    }>;
    type AuditFileTransferList = Common.PaginatingQueryRecord<AuditFileTransferLog>;
    type AuditFileTransferLogSearchParams = CommonType.RecordNullable<
      Pick<Api.Audit.AuditFileTransferLog, 'rustdesk_id' | 'peer_id' | 'uuid' | 'type' | 'created_at'> &
        Api.Common.CommonSearchParams
    >;
  }

  namespace System {
    type MailTemplate = Common.CommonRecord<{
      name: string;
      type: number;
      subject: string;
      contents: string;
    }>;
    type MailTemplateList = Common.PaginatingQueryRecord<MailTemplate>;
    type MailTemplateSearchParams = CommonType.RecordNullable<
      Pick<Api.System.MailTemplate, 'name' | 'type' | 'subject' | 'contents' | 'created_at'> &
        Api.Common.CommonSearchParams
    >;

    type MailLog = Common.CommonRecord<{
      username: string;
      uuid: string;
      subject: string;
      from: string;
      to: string;
      status: number;
    }>;
    type MailLogList = Common.PaginatingQueryRecord<MailLog>;
    type MailLogSearchParams = CommonType.RecordNullable<
      Pick<Api.System.MailLog, 'username' | 'uuid' | 'subject' | 'from' | 'to' | 'status' | 'created_at'> &
        Api.Common.CommonSearchParams
    >;
  }

  namespace AddressBooks {
    type AddressBook = Common.CommonRecord<{
      user_id: number;
      guid: string;
      name: string;
      owner: string;
      note: string;
      rule: number;
      max_peer: number;
      shared: boolean;
      peer_count: number;
      updated_at: string;
      actions?: any; // Virtual column for table actions
    }>;
    type AddressBooksList = Common.PaginatingQueryRecord<AddressBook>;
    type AddressBookDetail = AddressBook;
    type AddressBookSearchParams = CommonType.RecordNullable<
      Pick<Api.AddressBooks.AddressBook, 'name' | 'owner'> & Api.Common.CommonSearchParams
    >;

    type Peer = Common.CommonRecord<{
      rustdesk_id: string;
      username: string;
      hostname: string;
      alias: string;
      platform: string;
      tags: string;
      is_online: boolean;
      last_seen_at: string;
      ip_address: string;
    }>;
    type PeersList = {
      total: number;
      records: Peer[];
    };
  }

  /**
   * namespace Route
   *
   * backend api module: "route"
   */
  namespace Route {
    type ElegantConstRoute = import('@elegant-router/types').ElegantConstRoute;

    interface MenuRoute extends ElegantConstRoute {
      id: string;
    }

    interface UserRoute {
      routes: MenuRoute[];
      home: import('@elegant-router/types').LastLevelRouteKey;
    }
  }

  /**
   * namespace DocHelp
   *
   * backend api module: "dochelp"
   */
  namespace DocHelp {
    // Category
    type Category = Common.CommonRecord<{
      name: string;
      icon?: string;
      order: number;
    }>;

    type CategoryEdit = {
      name: string;
      icon?: string;
      order: number;
    };

    type CategoryList = {
      data: {
        categories: Category[];
      };
    };

    // Article
    type Article = Common.CommonRecord<{
      category_id: number;
      category_name: string;
      title: string;
      content: string;
      tags: string;
      is_pinned: boolean;
      views: number;
      author_id: number;
      author_name: string;
      from_ticket: number;
      updated_at: string;
    }>;

    type ArticleParams = Common.CommonSearchParams & {
      category_id?: number;
      search?: string;
      pinned?: string;
    };

    type ArticleList = {
      data: Common.PaginatingQueryRecord<Article>;
    };

    type ArticleDetail = {
      data: Article;
    };

    type ArticleEdit = {
      category_id: number;
      title: string;
      content: string;
      tags: string[];
      is_pinned: boolean;
    };

    // Ticket
    type Ticket = Common.CommonRecord<{
      title: string;
      description: string;
      status: number; // 1=Open, 2=In Analysis, 3=Resolved
      priority: number; // 1=Low, 2=Medium, 3=High, 4=Critical
      category_id: number;
      attachments: string;
      creator_id: number;
      creator_name: string;
      assigned_to: number;
      resolved_by: number;
      resolved_at: string;
      comment_count: number;
      updated_at: string;
    }>;

    type TicketParams = Common.CommonSearchParams & {
      status?: number;
      priority?: number;
      creator_id?: number;
    };

    type TicketList = {
      data: Common.PaginatingQueryRecord<Ticket>;
    };

    type TicketDetail = {
      data: Ticket;
    };

    type TicketCreate = {
      title: string;
      description: string;
      priority: number;
      category_id?: number;
      attachments?: string[];
    };

    type TicketUpdate = {
      status?: number;
      priority?: number;
      assigned_to?: number;
    };

    type TicketConvert = {
      category_id: number;
      title: string;
      tags: string[];
    };

    // Comment
    type Comment = Common.CommonRecord<{
      ticket_id: number;
      user_id: number;
      username: string;
      comment: string;
      attachments: string;
      is_internal: boolean;
    }>;

    type CommentList = {
      data: {
        comments: Comment[];
      };
    };

    type CommentCreate = {
      comment: string;
      attachments?: string[];
      is_internal?: boolean;
    };

    // Statistics
    type Stats = {
      data: {
        total_articles: number;
        open_tickets: number;
        in_progress_tickets: number;
        resolved_tickets: number;
        most_viewed_articles: Article[];
      };
    };
  }
}

