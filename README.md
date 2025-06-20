# uwud
> ultra whack unecessary dashboard


<p align="center">
  <img width=256 src="https://i.imgur.com/HepFtu3.png" alt="bubu heart" />
</p>

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/stefanistkuhl/uwud)
![Docker Image Size (tag)](https://img.shields.io/docker/image-size/raketenkult/uwud/latest)
### Why?
I wanted to change my homelab dashboard to be pink and pretty :3 But instead of either theming [Homer](https://github.com/bastienwirtz/homer) or looking for a different dashboard, I decided to make my own! It only has exactly the features I need with the look I want. This is probably the wrong dashboard for most people who want features besides a pretty pink aesthetic. Theming options to change the colors via the config will maybe be an option whenever I'm not too lazy to make it.

### Showcase

<details>
<summary>Dark Mode</summary>

![Dark mode](https://i.imgur.com/fRJ5Sa8.png)

</details>
<details>
<summary>Light Mode</summary>

![Light mode](https://i.imgur.com/PvXxbiH.png)

</details>

### Setup
To run `uwud`, create a directory and place a `config.toml` file within it, or utilize the [example](https://github.com/Stefanistkuhl/uwud/blob/master/config/config.example.toml) and map it to `/app/config`. Additionally, to use local image files for icons instead of links, create a directory and map it to `/app/static/images/custom`. Use either the `docker run` command or the Docker Compose configuration provided below, and modify values such as ports and restart policies as needed.
```sh
docker run -d \
  --name uwud \
  --restart unless-stopped \
  -v ./config:/app/config \
  -v ./images:/app/static/images/custom \
  -p 80:8080 \
  raketenkult/uwud:latest
```

```yml
services:
  uwud:
    container_name: uwud
    image: raketenkult/uwud:latest
    restart:  unless-stopped
    volumes:
      - ./config:/app/config
      - ./images:/app/static/images/custom
    ports:
      - 80:8080
```

### Configuration

The following configuration methods are available; there is also an [example](https://github.com/Stefanistkuhl/uwud/blob/master/config/config.example.toml) in the repository.

The `[general]` section contains `title`, `desc`, `icon`, and `colnums` options.

*   `title`: Used for the title in the header of the page and the page title itself. This is a non-optional value.
*   `desc`: The description displayed below the header. This is optional.
*   `icon`: Optional. This can be either a link to an image or the filename of an image. If you have bound a directory to `/app/static/images/custom`, simply provide the filename within your directory, as `/app/static/images/custom/` is automatically prepended to your configuration once loaded.
*   `colnums`: Required. This determines how many sections are displayed per row.

```toml
[general]
title="My Dashboard"
desc="A simple overview of my services."
icon="" # Can be any image URL or filename
columns=2
```

Sections are created using an array called `sections`, which have `name`, and `icon` values that function the same to those in the `[general]` section.

```toml
[[sections]]
name="Section One"
icon="" # Can be any image URL or filename
```

To add items to a section, `sections.items` is used. For each item, two new options are available: `url` and `health`. The `health` option is optional. URLs for health checks with API tokens are not supported. The health check is a simple HTTP GET request to the given URL, sent by the server running the dashboard, and only checks for a return code within the 200 range.

```toml
[[sections.items]]
name="Service A"
url="https://example.com/service-a" # Placeholder URL
desc="Description for Service A."
health="https://example.com/service-a/health" # Placeholder URL
icon="" # Can be any image URL or filename
```

To add more sections, simply add a new `[[section]]` with as many `[[sections.items]]` as needed. This can be repeated for any desired number of items and sections.
