export default {
  logo: <span>Crestic</span>,
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
