<script>
  export let name
  export let width = 24
  export let height = 24
  export let fill = 'none'
  export let stroke = 'currentColor'
</script>

<svg 
  width={width}
  height={height}
  fill={fill}
  stroke={stroke}
  stroke-width="2"
  stroke-linecap="round"
  stroke-linejoin="round"
>
  <use xlink:href="/feather-sprite.svg#{name}"/>
</svg>