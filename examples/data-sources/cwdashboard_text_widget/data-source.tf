data "cwdashboard_text_widget" "this" {
  markdown   = "# Hello, World!"
  background = "#000000"
  width      = 24
  height     = 2
}

data "cwdashboard" "this" {
  start           = "-PT7D"
  period_override = "auto"
  widgets = [
    data.cwdashboard_text_widget.this.json,
  ]
}
