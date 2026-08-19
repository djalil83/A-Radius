/* global window */
(function () {
  'use strict';

  const serviceToApi = {
    FTTH: 'FTTH',
    PPPoE: 'PPPOE',
    Voucher: 'HOTSPOT_VOUCHER',
    'Static IP': 'STATIC_IP'
  };
  const serviceFromApi = Object.fromEntries(
    Object.entries(serviceToApi).map(([label, value]) => [value, label])
  );

  function bps(value) {
    if (value === null || value === undefined || value === '') return null;
    const match = String(value).trim().match(/^(\d+(?:\.\d+)?)\s*(k|m|g)?\s*(?:bps)?$/i);
    if (!match) return null;
    const multiplier = { k: 1000, m: 1000000, g: 1000000000 }[String(match[2] || '').toLowerCase()] || 1;
    return Math.round(Number(match[1]) * multiplier);
  }

  function speed(value) {
    if (!value) return '';
    const n = Number(value);
    if (!Number.isFinite(n)) return '';
    if (n >= 1000000000) return `${n / 1000000000} Gbps`;
    if (n >= 1000000) return `${n / 1000000} Mbps`;
    if (n >= 1000) return `${n / 1000} Kbps`;
    return `${n} bps`;
  }

  function toUi(profile) {
    return {
      ...profile,
      id: String(profile.id),
      version: Number(profile.version || 1),
      service: serviceFromApi[profile.service_type] || profile.service_type,
      color: profile.color || '#1677ff',
      rate: profile.rate_limit || '',
      price: Number(profile.monthly_price || 0),
      status: profile.status === 'ACTIVE',
      mt: profile.mikrotik_group || '',
      radius: profile.radius_group || '',
      up: speed(profile.upload_bps),
      down: speed(profile.download_bps),
      shared: Number(profile.shared_users || 1)
    };
  }

  function toApi(profile) {
    const autoIsolate = Boolean(profile.isolate ?? true);
    return {
      name: profile.name,
      service_type: serviceToApi[profile.service] || profile.service,
      color: profile.color,
      rate_limit: profile.rate || null,
      upload_bps: bps(profile.up),
      download_bps: bps(profile.down),
      shared_users: Number(profile.shared || 1),
      monthly_price: Number(profile.price || 0),
      mikrotik_group: profile.mt || null,
      radius_group: profile.radius || null,
      auto_isolate: autoIsolate,
      version: Number(profile.version || 1)
    };
  }

  class ProfileAPIError extends Error {
    constructor(status, code, message) {
      super(message || 'API request failed');
      this.name = 'ProfileAPIError';
      this.status = status;
      this.code = code || 'HTTP_ERROR';
    }
  }

  class ProfileAPI {
    constructor(config = {}) {
      this.baseUrl = (config.baseUrl || window.ARADIUS_API_BASE || '').replace(/\/$/, '');
      this.tenantId = config.tenantId || window.ARADIUS_TENANT_ID || '';
      this.actorId = config.actorId || window.ARADIUS_ACTOR_ID || '';
    }

    headers() {
      if (!this.tenantId || !this.actorId) {
        throw new ProfileAPIError(401, 'IDENTITY_NOT_CONFIGURED', 'Tenant dan actor belum dikonfigurasi.');
      }
      return {
        Accept: 'application/json',
        'Content-Type': 'application/json',
        'X-Tenant-ID': this.tenantId,
        'X-Actor-ID': this.actorId
      };
    }

    async request(path, options = {}) {
      let response;
      try {
        response = await fetch(`${this.baseUrl}${path}`, {
          credentials: 'same-origin',
          ...options,
          headers: { ...this.headers(), ...(options.headers || {}) }
        });
      } catch (error) {
        throw new ProfileAPIError(0, 'NETWORK_ERROR', 'API tidak dapat dihubungi.');
      }
      const body = response.status === 204 ? null : await response.json().catch(() => null);
      if (!response.ok) {
        const detail = body && body.error;
        throw new ProfileAPIError(response.status, detail && detail.code, detail && detail.message);
      }
      return body;
    }

    async list(filters = {}) {
      const query = new URLSearchParams({ limit: '100', offset: '0' });
      if (filters.q) query.set('q', filters.q);
      if (filters.service_type) query.set('service_type', filters.service_type);
      if (filters.status) query.set('status', filters.status);
      const result = await this.request(`/api/v1/subscription-profiles?${query}`);
      return (result.items || []).map(toUi);
    }

    async create(profile) {
      return toUi(await this.request('/api/v1/subscription-profiles', {
        method: 'POST', body: JSON.stringify(toApi(profile))
      }));
    }

    async update(profile) {
      return toUi(await this.request(`/api/v1/subscription-profiles/${encodeURIComponent(profile.id)}`, {
        method: 'PATCH', body: JSON.stringify(toApi(profile))
      }));
    }

    async archive(profile) {
      await this.request(`/api/v1/subscription-profiles/${encodeURIComponent(profile.id)}?version=${encodeURIComponent(profile.version)}`, {
        method: 'DELETE'
      });
    }

    async revisions(profile) {
      return this.request(`/api/v1/subscription-profiles/${encodeURIComponent(profile.id)}/revisions?limit=100`);
    }
  }

  window.ProfileAPI = ProfileAPI;
  window.ProfileAPIError = ProfileAPIError;
})();
