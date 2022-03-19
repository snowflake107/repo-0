<div id="top"></div>
<!--
*** Thanks for checking out the Best-README-Template. If you have a suggestion
*** that would make this better, please fork the repo and create a pull request
*** or simply open an issue with the tag "enhancement".
*** Don't forget to give the project a star!
*** Thanks again! Now go create something AMAZING! :D
-->

<!-- PROJECT SHIELDS -->
<!--
*** I'm using markdown "reference style" links for readability.
*** Reference links are enclosed in brackets [ ] instead of parentheses ( ).
*** See the bottom of this document for the declaration of the reference variables
*** for contributors-url, forks-url, etc. This is an optional, concise syntax you may use.
*** https://www.markdownguide.org/basic-syntax/#reference-style-links
-->

[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![MIT License][license-shield]][license-url]
[![LinkedIn][linkedin-shield]][linkedin-url]

<!-- PROJECT LOGO -->
<br />
<div align="center">
  <a href="https://www.mautic.org/">
    <img src="https://www.mautic.org/themes/custom/mauticorg_base/logo.svg" alt="Logo" width="425" height="113">
  </a>

  <h3 align="center">mautic</h3>

  <p align="center">
    An unofficial docker container for Mautic that works
    <br />
    <a href="https://github.com/aperim/docker-mautic"><strong>Explore the docs »</strong></a>
    <br />
    <br />
    <a href="https://www.mautic.org/what-is-mautic">What is mautic</a>
    ·
    <a href="https://docs.mautic.org/en">Docs</a>
    ·
    <a href="https://www.mautic.org/blog">Blog</a>
  </p>
</div>

<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#built-with">Built With</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#roadmap">Roadmap</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#license">License</a></li>
    <li><a href="#contact">Contact</a></li>
    <li><a href="#acknowledgments">Acknowledgments</a></li>
  </ol>
</details>

<!-- ABOUT THE PROJECT -->

## About The Project

[![Product Name Screen Shot][product-screenshot]](https://mautic.org)

We could never get the official mautic docker containers to work. They'd always fail if they ever worked in the first place - so we made a container that works, and can even keep your mautic version up-to-date.

This is the best docker container to use for mautic, don't waste your time trying to get others working.

Here's why:

- This container works
- It doesn't waste your time
- It works

Of course, there's no guarantee it will actually work for you. It does for us, so that's great - so it should work for you. No promises though.

Use the `docker-compose.yaml` to get started.

<p align="right">(<a href="#top">back to top</a>)</p>

### Built With

We aren't responsible for mautic - it's a great tool, and we love it. We are only responsible for getting it into a working container.

<p align="right">(<a href="#top">back to top</a>)</p>

<!-- GETTING STARTED -->

## Getting Started

Use the example compose file to see how you can use mautic in a container.

### Prerequisites

You need docker, a database server and other things needed for successful operation of a business, we can't go into that here.

### Installation

_There's not really anything to install per-se._

1. Get Docker [https://docker.com](https://docker.com)
2. Create a docker compose file
   ```sh
   vi docker-compose.yaml
   ```
3. Configure it as needed
4. Run it
   ```sh
   docker-compose up -d
   ```

<p align="right">(<a href="#top">back to top</a>)</p>

<!-- USAGE EXAMPLES -->

## Usage

There's only really one use for this container - running mautic.

_For more examples, please refer to the [Documentation](https://mautic.org)_

<p align="right">(<a href="#top">back to top</a>)</p>

<!-- ROADMAP -->

## Roadmap

- [x] Add Readme
- [x] Add example compose
- [ ] Make the readme actually helpful
- [ ] Improve the starup script
- [ ] Better examples
  - [ ] Full with mysql
  - [ ] Not so fancy

See the [open issues](https://github.com/aperim/docker-mautic/issues) for a full list of proposed features (and known issues).

<p align="right">(<a href="#top">back to top</a>)</p>

<!-- CONTRIBUTING -->

## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

<p align="right">(<a href="#top">back to top</a>)</p>

<!-- LICENSE -->

## License

This is distributed under the Apache License. See `LICENSE` for more information.

<p align="right">(<a href="#top">back to top</a>)</p>

<!-- CONTACT -->

## Contact

Your Name - [@aperimstudio](https://twitter.com/aperimstudio) - hello@aperim.com

Project Link: [https://github.com/mautic](https://github.com/mautic)

<p align="right">(<a href="#top">back to top</a>)</p>

<!-- ACKNOWLEDGMENTS -->

## Acknowledgments

This collection of files was created on the unceded Aboriginal land. We acknowledge the Gadigal people of the Eora Nation as the Traditional Custodians of the Country we work on. We recognise their continuing connection to the land and waters, and thank them for protecting this coastline and its ecosystems since time immemorial.

We'd like to thank the Internet, for existing. The mautic team and contributers for ... mautic. And you, for (hopefully) making the world a better place.

<p align="right">(<a href="#top">back to top</a>)</p>

<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->

[contributors-shield]: https://img.shields.io/github/contributors/aperim/docker-mautic.svg?style=for-the-badge
[contributors-url]: https://github.com/aperim/docker-mautic/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/aperim/docker-mautic.svg?style=for-the-badge
[forks-url]: https://github.com/aperim/docker-mautic/network/members
[stars-shield]: https://img.shields.io/github/stars/aperim/docker-mautic.svg?style=for-the-badge
[stars-url]: https://github.com/aperim/docker-mautic/stargazers
[issues-shield]: https://img.shields.io/github/issues/aperim/docker-mautic.svg?style=for-the-badge
[issues-url]: https://github.com/aperim/docker-mautic/issues
[license-shield]: https://img.shields.io/github/license/aperim/docker-mautic.svg?style=for-the-badge
[license-url]: https://github.com/aperim/docker-mautic/blob/master/LICENSE.txt
[linkedin-shield]: https://img.shields.io/badge/-LinkedIn-black.svg?style=for-the-badge&logo=linkedin&colorB=555
[linkedin-url]: https://www.linkedin.com/company/aperim
[product-screenshot]: images/mautic_screenshot.png
