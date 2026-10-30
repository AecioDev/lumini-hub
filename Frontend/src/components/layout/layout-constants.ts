// Altura do bloco do topo do menu lateral e do cabeçalho: iguais, para a linha divisória dos dois
// ficar alinhada de ponta a ponta.
export const TOP_BAR_HEIGHT = 64;

// Larguras (px) em que o menu muda de comportamento (as mesmas do `md` e do `lg` do antd):
//   >= BREAKPOINT_DESKTOP: menu lateral fixo, completo ou recolhido (preferência do usuário);
//   >= BREAKPOINT_TABLET e abaixo disso: menu lateral fixo só com ícones, e o botão do cabeçalho abre a gaveta;
//   abaixo de BREAKPOINT_TABLET: sem menu lateral fixo, só a gaveta.
export const BREAKPOINT_DESKTOP = 992;
export const BREAKPOINT_TABLET = 768;
