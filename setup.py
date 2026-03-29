"""
Setup script for auto-workflow.
"""
from setuptools import setup, find_packages

with open("README_PYTHON.md", "r", encoding="utf-8") as fh:
    long_description = fh.read()

setup(
    name="auto-workflow",
    version="1.0.0",
    author="Auto Workflow Team",
    description="A flexible automation workflow system",
    long_description=long_description,
    long_description_content_type="text/markdown",
    url="https://github.com/yourusername/auto-workflow",
    packages=find_packages(),
    classifiers=[
        "Development Status :: 4 - Beta",
        "Intended Audience :: Developers",
        "Topic :: Software Development :: Libraries :: Python Modules",
        "License :: OSI Approved :: MIT License",
        "Programming Language :: Python :: 3",
        "Programming Language :: Python :: 3.8",
        "Programming Language :: Python :: 3.9",
        "Programming Language :: Python :: 3.10",
        "Programming Language :: Python :: 3.11",
    ],
    python_requires=">=3.8",
    install_requires=[
        "requests>=2.31.0",
    ],
    entry_points={
        "console_scripts": [
            "auto-workflow=main:main",
        ],
    },
)
