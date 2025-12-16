<script lang="ts">
	import { SOCIETIES, type Society } from "$lib/types";
    import { Stateful } from "../shared.svelte";
    let person = $state(Stateful.person)


    function handleSociety(s: Society) {
      console.log("SET TO ", s)
      Stateful.person.Society = s
    }

</script>

<div class="flex flex-col items-center mt-10 gap-4">
  <div class="form-control">
    <label class="label" for="firstName">First Name</label>
    <input
      id="firstName"
      type="text"
      class="input input-bordered"
      placeholder="first name"
      required
      bind:value={person.FirstName}
    />
  </div>

  <div class="form-control">
    <label class="label" for="lastName">Last Name</label>
    <input
      id="lastName"
      type="text"
      class="input input-bordered"
      placeholder="last name"
      required
      bind:value={person.LastName}
    />
  </div>

  <div class="form-control">
    <label class="label" for="writerIpi">Writer IPI</label>
    <input
      id="writerIpi"
      type="number"
      class="input input-bordered"
      placeholder="writer ipi"
      required
      bind:value={person.WriterIpiNum}
    />
  </div>

  <div class="form-control">
    <label class="label" for="pubIpi">Publisher IPI</label>
    <input
      id="pubIpi"
      type="number"
      class="input input-bordered"
      placeholder="publisher ipi"
      required
      bind:value={person.PubIpiNum}
    />
  </div>

  <fieldset class="fieldset">
    <legend class="fieldset-legend">Society</legend>
    <select class="select">
      <option disabled selected>pick a society</option>
      {#each SOCIETIES as s}
        <option onclick={() => handleSociety(s)}>{s}</option>
      {/each}
    </select>
  </fieldset>
</div>
