package webnamesru

import (
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
)

func init() {
	caddy.RegisterModule(Provider{})
}

func (Provider) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "dns.providers.webnamesru",
		New: func() caddy.Module { return new(Provider) },
	}
}

func (p *Provider) Provision(ctx caddy.Context) error {
	return nil
}

// UnmarshalCaddyfile поддерживает синтаксис:
//   dns webnamesru {env.WEBNAMESRU_API_KEY}
// или
//   dns webnamesru {
//       api_key <ключ>
//   }
func (p *Provider) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	d.Next() // пропускаем имя директивы "webnamesru"
	if d.NextArg() {
		// Инлайн-форма: dns webnamesru <api_key>
		p.APIKey = d.Val()
		return nil
	}
	// Блочная форма
	for nesting := d.Nesting(); d.NextBlock(nesting); {
		switch d.Val() {
		case "api_key":
			if !d.NextArg() {
				return d.ArgErr()
			}
			p.APIKey = d.Val()
		default:
			return d.Errf("неизвестная опция: %s", d.Val())
		}
	}
	return nil
}
