export default {
  logo: (
    <>
      <img src="/logo.png" width={24} height={24} alt="Crestic" />
      <span style={{ marginLeft: '.4em', fontWeight: 600 }}>Crestic</span>
    </>
  ),
  head: (
    <>
      <link rel="icon" type="image/png" href="/favicon.png" />
      <link rel="apple-touch-icon" href="/apple-touch-icon.png" />
      <meta property="og:image" content="https://crestic.kolodka.fyi/og.png" />
      <meta name="twitter:card" content="summary" />
    </>
  ),
  docsRepositoryBase: 'https://github.com/alexander-kolodka/crestic/tree/main/docs',
  project: {
    link: 'https://github.com/alexander-kolodka/crestic',
  },
  sidebar: {
    defaultMenuCollapseLevel: 1,
  },
  feedback: {
    content: 'Question? An error? Give feedback →',
  },
  footer: {
    text: (
      <span>
        MIT {new Date().getFullYear()} ©{' '}
        <a href="https://github.com/alexander-kolodka" target="_blank">
          Alexander Kolodka
        </a>
      </span>
    ),
  },
  useNextSeoProps() {
    return {
      titleTemplate: '%s – Crestic',
      defaultTitle: 'Crestic',
    };
  },
}
