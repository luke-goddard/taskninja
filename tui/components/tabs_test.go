package components

import (
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tabs", func() {

	var tabs *Tabs

	BeforeEach(func() {
		tabs = NewTabs()
	})
	It("should have 8 tabs", func() {
		Expect(len(tabs.Tabs)).To(Equal(8))
	})
	It("should have the first tab as 'Tasks'", func() {
		Expect(tabs.Tabs[0]).To(Equal("Tasks (1)"))
	})
	It("should have the second tab as 'Projects'", func() {
		Expect(tabs.Tabs[1]).To(Equal("Projects (2)"))
	})
	It("should have the third tab as 'Tags'", func() {
		Expect(tabs.Tabs[2]).To(Equal("Tags (3)"))
	})
	It("should have the fourth tab as 'People'", func() {
		Expect(tabs.Tabs[3]).To(Equal("People (4)"))
	})
	It("should have the first tab as active", func() {
		Expect(tabs.ActiveTab).To(Equal(0))
	})
	It("should move RIGHT when pressing 'l'", func() {
		Expect(tabs.ActiveTab).To(Equal(0))
		tabs, _ = tabs.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})
		Expect(tabs.ActiveTab).To(Equal(1))
	})
	It("should move RIGHT when pressing 'RightArrow'", func() {
		Expect(tabs.ActiveTab).To(Equal(0))
		tabs, _ = tabs.Update(tea.KeyMsg{Type: tea.KeyRight})
		Expect(tabs.ActiveTab).To(Equal(1))
	})
	It("should move RIGHT when pressing 'n'", func() {
		Expect(tabs.ActiveTab).To(Equal(0))
		tabs, _ = tabs.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
		Expect(tabs.ActiveTab).To(Equal(1))
	})
	It("should move RIGHT when pressing 'tab'", func() {
		Expect(tabs.ActiveTab).To(Equal(0))
		tabs, _ = tabs.Update(tea.KeyMsg{Type: tea.KeyTab})
		Expect(tabs.ActiveTab).To(Equal(1))
	})
	It("should move LEFT when pressing 'h'", func() {
		Expect(tabs.ActiveTab).To(Equal(0))
		tabs, _ = tabs.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})
		tabs, _ = tabs.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})
		Expect(tabs.ActiveTab).To(Equal(0))
	})
	It("should move LEFT when pressing 'LeftArrow'", func() {
		Expect(tabs.ActiveTab).To(Equal(0))
		tabs, _ = tabs.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})
		tabs, _ = tabs.Update(tea.KeyMsg{Type: tea.KeyLeft})
		Expect(tabs.ActiveTab).To(Equal(0))
	})
	It("should move LEFT when pressing 'p'", func() {
		Expect(tabs.ActiveTab).To(Equal(0))
		tabs, _ = tabs.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})
		tabs, _ = tabs.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}})
		Expect(tabs.ActiveTab).To(Equal(0))
	})
	It("should move LEFT when pressing 'shift+tab'", func() {
		Expect(tabs.ActiveTab).To(Equal(0))
		tabs, _ = tabs.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})
		tabs, _ = tabs.Update(tea.KeyMsg{Type: tea.KeyShiftTab})
		Expect(tabs.ActiveTab).To(Equal(0))
	})

	Describe("When the active tab is the last tab", func() {
		BeforeEach(func() {
			tabs.ActiveTab = 7
		})

		It("It should not change tabs when moving right", func() {
			Expect(tabs.ActiveTab).To(Equal(7))
			tabs, _ = tabs.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})
			Expect(tabs.ActiveTab).To(Equal(7))
		})
	})

	Describe("When the active tab is the first tab", func() {
		BeforeEach(func() {
			tabs.ActiveTab = 0
		})

		It("It should not change tabs when moving left", func() {
			Expect(tabs.ActiveTab).To(Equal(0))
			tabs, _ = tabs.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})
			Expect(tabs.ActiveTab).To(Equal(0))
		})
	})
})
