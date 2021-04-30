<template>
  <div class="container">
    <div class="footer">
      <div class="backdrop"></div>
      <div class="items">
        <div class="left">
          <div class="cta">
            Sign up for Sifchain updates<input /><button>Stay Informed</button>
          </div>
        </div>
        <div class="right">
          <div class="links">
            <a href="">Privacy Policy</a>
            <a href="">Roadmap</a>
            <a href="">Legal Disclaimer</a>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.container {
  position: fixed;
  width: 100vw;
  height: 51px;
  bottom: 0;
  left: 0;
}
.footer {
  font-size: 12px; /* make $variable */
  line-height: 14px; /* make $variable */
  text-align: left;
  font-weight: 400; /* make $variable */
  position: relative;
}
.items {
  display: flex;
  position: absolute;
  width: 100%;
}

.left {
  flex: 1 1 auto;
  height: 51px;
}
.right {
  flex: 1 1 auto;
  height: 51px;
}
.cta {
  display: flex;
  align-items: center;
  padding-left: 15px;
  height: 51px;
  color: #fff;
}
.links {
  display: flex;
  justify-content: flex-end;
  padding-right: 15px;
  height: 51px;
  align-items: center;
  a {
    margin-left: 15px;
    color: #fff;
  }
}
.backdrop {
  background: black;
  position: absolute;
  width: 100%;
  opacity: 0.5;
  height: 51px;
}
</style>
