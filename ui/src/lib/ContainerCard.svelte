<script lang="ts">
  import type { Container } from "./api";

  export let container: Container;
</script>

<a
  href="{container.targetUrl}" target="_blank" rel="noreferrer"
  class="block px-5 pt-5 m-5 flex flex-col items-center text-center dark:bg-slate-700 bg-white shadow-md dark:shadow-slate-700 shadow-neutral-100 dark:border-slate-600 border dark:hover:bg-slate-600 hover:bg-slate-50 cursor-pointer"
>
  <p class="font-semibold text-xl  dark:text-white text-slate-700">{container.service}</p>

  <p class="pb-2 pt-5 font-light text-sm  dark:text-slate-500 text-gray-400">
    running {container.image}
  </p>
</a>
