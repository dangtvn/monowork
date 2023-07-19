export type SiteConfig = typeof siteConfig

export const siteConfig = {
  name: "Mono work",
  description:
    "Beautifully designed components built with Radix UI and Tailwind CSS.",
  mainNav: [
    {
      title: "Home",
      href: "/",
    },
    {
      title: "Login",
      href: "/login",
    },
  ],
  links: {
    twitter: "https://twitter.com/dangtvn",
    github: "https://github.com/dangtvn",
  },
}
